package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {
	cfg := LoadUploadConfig()
	files, err := ParseUploads(r, "files", cfg)
	if err != nil {
		uploadErrorResponse(w, err, cfg)
		return
	}
	if len(files) != 1 {
		WriteError(w, "kirim tepat 1 file", http.StatusBadRequest)
		return
	}

	dir, err := NewWorkDir("watermark-*")
	if err != nil {
		WriteError(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	in, err := saveUpload(files[0], dir, "watermark", -1, []uploadedKind{pdfKind}, "PDF")
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := filepath.Join(dir, "watermarked.pdf")
	var processErr error
	marks := r.MultipartForm.File["mark"]
	switch len(marks) {
	case 0:
		processErr = pdf.Watermark(r.Context(), in, out,
			r.FormValue("text"), r.FormValue("position"), r.FormValue("opacity"), r.FormValue("size"))
	case 1:
		img, err := saveUpload(marks[0], dir, "mark", -1, []uploadedKind{jpegKind, pngKind}, "gambar PNG/JPEG")
		if err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		processErr = pdf.WatermarkImage(r.Context(), in, out, img,
			r.FormValue("position"), r.FormValue("opacity"), r.FormValue("size"))
	default:
		WriteError(w, "kirim maksimal 1 gambar watermark", http.StatusBadRequest)
		return
	}

	if processErr != nil {
		if errors.Is(processErr, pdf.ErrParam) {
			WriteError(w, processErr.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("pemrosesan watermark gagal", "err", processErr)
		WriteError(w, "watermark tidak bisa ditambahkan", http.StatusUnprocessableEntity)
		return
	}
	if err := WriteFileResponse(w, out, "watermarked.pdf", "application/pdf"); err != nil {
		WriteError(w, "gagal mengirim hasil", http.StatusInternalServerError)
	}
}