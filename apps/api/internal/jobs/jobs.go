// Package jobs mengelola antrian pekerjaan asynchronous (asynq + Redis).
// Tidak ada database: metadata disimpan di Redis dengan TTL, file di folder
// sementara per job yang dibersihkan oleh sweeper.
package jobs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// TopicOffice adalah nama tugas asynq untuk konversi Office ke PDF.
const TopicOffice = "office-to-pdf"

// TTL menentukan berapa lama metadata Redis dan folder job disimpan.
const TTL = time.Hour

// ErrNotFound menandakan job tidak ada atau sudah kedaluwarsa.
var ErrNotFound = errors.New("job tidak ditemukan")

var idPattern = regexp.MustCompile(`^[a-f0-9]{32}$`)

// ValidID memeriksa format id job agar aman dipakai sebagai bagian path.
func ValidID(id string) bool {
	return idPattern.MatchString(id)
}

func redisAddr() string {
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		return v
	}
	return "localhost:6379"
}

// AsynqOpt mengembalikan koneksi Redis untuk asynq.
func AsynqOpt() asynq.RedisClientOpt {
	return asynq.RedisClientOpt{Addr: redisAddr()}
}

func store() *redis.Client {
	return redis.NewClient(&redis.Options{Addr: redisAddr()})
}

// Payload adalah argumen tugas asynq. Semua path milik server, bukan input user.
type Payload struct {
	JobID  string `json:"job_id"`
	Tool   string `json:"tool"`
	Input  string `json:"input"`
	OutDir string `json:"outdir"`
}

// Meta adalah status job yang disimpan di Redis dengan TTL.
type Meta struct {
	Tool   string `json:"tool"`
	Status string `json:"status"` // queued | processing | done | error
	Error  string `json:"error,omitempty"`
}

func metaKey(id string) string {
	return "gidocs:job:" + id
}

// NewID membuat id acak 32 karakter heksadesimal.
func NewID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand tidak pernah gagal pada sistem normal
	}
	return hex.EncodeToString(b)
}

// JobDir mengembalikan folder kerja untuk satu job.
func JobDir(id string) string {
	return filepath.Join(os.TempDir(), "gidocs-jobs", id)
}

// ResultPath adalah lokasi PDF hasil setelah job selesai.
func ResultPath(id string) string {
	return filepath.Join(JobDir(id), "result.pdf")
}

// ZipResultPath adalah lokasi hasil berbentuk ZIP setelah job selesai
// (mis. pdf-to-image).
func ZipResultPath(id string) string {
	return filepath.Join(JobDir(id), "result.zip")
}

// Enqueue menyimpan meta awal lalu memasukkan tugas ke antrian Redis.
// Folder job harus sudah dibuat dan berisi file input.
func Enqueue(ctx context.Context, p Payload) error {
	if !ValidID(p.JobID) {
		return fmt.Errorf("id job tidak valid")
	}
	meta, err := json.Marshal(Meta{Tool: p.Tool, Status: "queued"})
	if err != nil {
		return err
	}
	c := store()
	if err := c.Set(ctx, metaKey(p.JobID), meta, TTL).Err(); err != nil {
		return err
	}
	payload, err := json.Marshal(p)
	if err != nil {
		return err
	}
	cli := asynq.NewClient(AsynqOpt())
	defer cli.Close()
	_, err = cli.EnqueueContext(ctx, asynq.NewTask(TopicOffice, payload),
		asynq.TaskID(p.JobID),
		asynq.Retention(TTL),
		asynq.MaxRetry(0))
	return err
}

// SetStatus memperbarui meta job (dipanggil worker) dan menyegarkan TTL.
func SetStatus(ctx context.Context, id, tool, status, errMsg string) error {
	meta, err := json.Marshal(Meta{Tool: tool, Status: status, Error: errMsg})
	if err != nil {
		return err
	}
	return store().Set(ctx, metaKey(id), meta, TTL).Err()
}

// GetMeta mengambil meta job; ErrNotFound jika tidak ada/kedaluwarsa.
func GetMeta(ctx context.Context, id string) (Meta, error) {
	if !ValidID(id) {
		return Meta{}, ErrNotFound
	}
	raw, err := store().Get(ctx, metaKey(id)).Result()
	if errors.Is(err, redis.Nil) {
		return Meta{}, ErrNotFound
	}
	if err != nil {
		return Meta{}, err
	}
	var m Meta
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return Meta{}, err
	}
	return m, nil
}

// Sweep menghapus folder job yang lebih tua dari ttl. Dipanggil berkala oleh worker.
func Sweep(root string, ttl time.Duration) (int, error) {
	entries, err := os.ReadDir(root)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, e := range entries {
		if !e.IsDir() || !ValidID(e.Name()) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if time.Since(info.ModTime()) > ttl {
			if err := os.RemoveAll(filepath.Join(root, e.Name())); err == nil {
				removed++
			}
		}
	}
	return removed, nil
}

// Root adalah folder induk semua job.
func Root() string {
	return filepath.Join(os.TempDir(), "gidocs-jobs")
}
