package handler

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const (
	defaultMaxFiles      = 20
	defaultMaxTotalBytes = 50 << 20 // 50 MB
)

// UploadConfig holds upload limits from env
type UploadConfig struct {
	MaxFiles      int
	MaxTotalBytes int64
}

// LoadUploadConfig reads upload limits from environment variables
func LoadUploadConfig() UploadConfig {
	maxFiles := defaultMaxFiles
	if v := os.Getenv("MAX_FILES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxFiles = parsed
		}
	}
	maxBytes := int64(defaultMaxTotalBytes)
	if v := os.Getenv("MAX_UPLOAD_MB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxBytes = int64(parsed) << 20
		}
	}
	return UploadConfig{MaxFiles: maxFiles, MaxTotalBytes: maxBytes}
}

// ParseUploads parses multipart form, validates file count, returns FileHeaders
func ParseUploads(r *http.Request, fieldName string, cfg UploadConfig) ([]*multipart.FileHeader, error) {
	r.Body = http.MaxBytesReader(nil, r.Body, cfg.MaxTotalBytes)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, ErrUploadTooLarge
	}
	files := r.MultipartForm.File[fieldName]
	if len(files) == 0 {
		return nil, ErrNoFiles
	}
	if len(files) > cfg.MaxFiles {
		return nil, ErrTooManyFiles
	}
	return files, nil
}

// SavePDF validates magic bytes "%PDF-" and saves with generated name.
// index < 0 membuat nama tunggal (base saja), selain itu base-N.
func SavePDF(fh *multipart.FileHeader, dir, base string, index int) (string, error) {
	return saveUpload(fh, dir, base, index, []uploadedKind{pdfKind}, "PDF")
}

// SaveImage validates magic bytes JPEG/PNG and saves with generated name.
func SaveImage(fh *multipart.FileHeader, dir, base string, index int) (string, error) {
	return saveUpload(fh, dir, base, index, []uploadedKind{jpegKind, pngKind}, "gambar JPEG/PNG")
}

type uploadedKind struct {
	magic []byte
	ext   string
}

var (
	pdfKind  = uploadedKind{[]byte("%PDF-"), ".pdf"}
	jpegKind = uploadedKind{[]byte{0xFF, 0xD8, 0xFF}, ".jpg"}
	pngKind  = uploadedKind{[]byte{0x89, 'P', 'N', 'G'}, ".png"}
)

func saveUpload(fh *multipart.FileHeader, dir, base string, index int, kinds []uploadedKind, label string) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", errors.New("gagal membuka file upload")
	}
	defer src.Close()

	head := make([]byte, 8)
	if _, err := io.ReadFull(src, head); err != nil {
		return "", badUploadError(index, label)
	}
	ext := ""
	for _, k := range kinds {
		if len(head) >= len(k.magic) && bytes.Equal(head[:len(k.magic)], k.magic) {
			ext = k.ext
			break
		}
	}
	if ext == "" {
		return "", badUploadError(index, label)
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("gagal membaca file upload")
	}

	name := base
	if index >= 0 {
		name += "-" + strconv.Itoa(index)
	}
	path := filepath.Join(dir, name+ext)
	dst, err := os.Create(path)
	if err != nil {
		return "", errors.New("gagal menyimpan file sementara")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("gagal menyimpan file sementara")
	}
	return path, nil
}

// NewWorkDir creates a temp directory for request-scoped work
func NewWorkDir(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// WriteError writes a JSON/text error response
func WriteError(w http.ResponseWriter, msg string, code int) {
	http.Error(w, msg, code)
}

// WriteFileResponse streams a file as attachment
func WriteFileResponse(w http.ResponseWriter, path, filename, contentType string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	_, err = io.Copy(w, f)
	return err
}

// Sentinel errors for caller to handle
var (
	ErrUploadTooLarge = errors.New("upload terlalu besar atau tidak valid")
	ErrNoFiles        = errors.New("tidak ada file di field 'files'")
	ErrTooManyFiles   = errors.New("terlalu banyak file (maks 20)")
)

func badUploadError(index int, label string) error {
	i := index + 1
	if index < 0 {
		i = 1
	}
	return errors.New("file ke-" + strconv.Itoa(i) + " bukan " + label + " yang valid")
}