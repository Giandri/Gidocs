package http

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

const (
	maxFiles      = 20
	maxTotalBytes = 50 << 20 // 50 MB
)

func MergeHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxTotalBytes)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "upload terlalu besar atau tidak valid", http.StatusRequestEntityTooLarge)
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) < 2 || len(files) > maxFiles {
		http.Error(w, "kirim 2 sampai 20 file PDF pada field 'files'", http.StatusBadRequest)
		return
	}

	dir, err := os.MkdirTemp("", "merge-*")
	if err != nil {
		http.Error(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	inputs := make([]string, 0, len(files))
	for i, fh := range files {
		path, err := savePDF(fh, dir, i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		inputs = append(inputs, path)
	}

	output := filepath.Join(dir, "merged.pdf")
	if err := pdf.Merge(r.Context(), inputs, output); err != nil {
		http.Error(w, "gagal menggabungkan PDF (file mungkin rusak atau terkunci)", http.StatusUnprocessableEntity)
		return
	}

	f, err := os.Open(output)
	if err != nil {
		http.Error(w, "gagal membaca hasil", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="merged.pdf"`)
	io.Copy(w, f)
}

// savePDF memeriksa magic bytes "%PDF-" lalu menyimpan dengan nama buatan sendiri.
func savePDF(fh *multipart.FileHeader, dir string, i int) (string, error) {
	src, err := fh.Open()
	if err != nil {
		return "", errors.New("gagal membuka file upload")
	}
	defer src.Close()

	head := make([]byte, 5)
	if _, err := io.ReadFull(src, head); err != nil || !bytes.Equal(head, []byte("%PDF-")) {
		return "", errors.New("file ke-" + strconv.Itoa(i+1) + " bukan PDF yang valid")
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", errors.New("gagal membaca file upload")
	}

	path := filepath.Join(dir, "in-"+strconv.Itoa(i)+".pdf")
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