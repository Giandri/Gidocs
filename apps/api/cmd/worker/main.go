// Worker menjalankan antrian tugas asynq: konversi Office ke PDF via
// LibreOffice headless, plus pembersih folder job kedaluwarsa.
package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Giandri/Gidocs/apps/api/internal/jobs"
	"github.com/hibiken/asynq"
)

func main() {
	go sweeper()

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TopicOffice, handleOffice)

	// Serial (1 job pada satu waktu): seluruh job berbagi satu profil
	// LibreOffice, jadi dua konversi tidak boleh berjalan bersamaan.
	srv := asynq.NewServer(jobs.AsynqOpt(), asynq.Config{Concurrency: 1})
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		srv.Shutdown()
	}()

	slog.Info("worker berjalan", "redis", "localhost:6379")
	if err := srv.Run(mux); err != nil {
		slog.Error("worker gagal berjalan", "err", err)
		os.Exit(1)
	}
}

// handleOffice memproses tugas konversi Office. Selalu mengembalikan nil
// setelah meta diupdate agar asynq tidak mengulang job yang sudah diproses.
func handleOffice(ctx context.Context, t *asynq.Task) error {
	var p jobs.Payload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		slog.Error("payload job tidak valid", "err", err)
		return nil
	}
	if !jobs.ValidID(p.JobID) {
		slog.Error("id job tidak valid")
		return nil
	}

	jobs.SetStatus(ctx, p.JobID, p.Tool, "processing", "")

	if err := jobs.ConvertOffice(ctx, p); err != nil {
		// pesan error aman ditampilkan ke user (berbahasa Indonesia, generik)
		jobs.SetStatus(ctx, p.JobID, p.Tool, "error", err.Error())
		return nil
	}
	jobs.SetStatus(ctx, p.JobID, p.Tool, "done", "")
	return nil
}

// sweeper menghapus folder job yang lebih tua dari TTL secara berkala.
func sweeper() {
	for {
		time.Sleep(10 * time.Minute)
		removed, err := jobs.Sweep(jobs.Root(), jobs.TTL)
		if err != nil {
			slog.Error("pembersih job gagal", "err", err)
			continue
		}
		if removed > 0 {
			slog.Info("folder job dibersihkan", "jumlah", removed)
		}
	}
}
