package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func MergeHandler(w http.ResponseWriter, r *http.Request) {
	cfg := LoadUploadConfig()

	files, err := ParseUploads(r, "files", cfg)
	if err != nil {
		switch err {
		case ErrUploadTooLarge:
			WriteError(w, "upload terlalu besar atau tidak valid", http.StatusRequestEntityTooLarge)
		case ErrNoFiles:
			WriteError(w, "kirim 2 sampai 20 file PDF pada field 'files'", http.StatusBadRequest)
		case ErrTooManyFiles:
			WriteError(w, "terlalu banyak file (maks "+strconv.Itoa(cfg.MaxFiles)+")", http.StatusBadRequest)
		default:
			WriteError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if len(files) < 2 {
		WriteError(w, "kirim 2 sampai 20 file PDF pada field 'files'", http.StatusBadRequest)
		return
	}

	dir, err := NewWorkDir("merge-*")
	if err != nil {
		WriteError(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	inputs := make([]string, 0, len(files))
	for i, fh := range files {
		path, err := SavePDF(fh, dir, "in", i)
		if err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		inputs = append(inputs, path)
	}

	output := filepath.Join(dir, "merged.pdf")
	if err := pdf.Merge(r.Context(), inputs, output); err != nil {
		WriteError(w, "gagal menggabungkan PDF (file mungkin rusak atau terkunci)", http.StatusUnprocessableEntity)
		return
	}

	if err := WriteFileResponse(w, output, "merged.pdf", "application/pdf"); err != nil {
		WriteError(w, "gagal mengirim hasil", http.StatusInternalServerError)
	}
}