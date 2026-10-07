package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

// SignHandler menempelkan gambar tanda tangan ke satu halaman PDF.
// Input: 1 PDF di field "files", 1 gambar PNG/JPEG di field "mark",
// posisi di field page, x, y, width (poin PDF).
func SignHandler(w http.ResponseWriter, r *http.Request) {
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

	marks := r.MultipartForm.File["mark"]
	if len(marks) != 1 {
		WriteError(w, "kirim tepat 1 gambar tanda tangan", http.StatusBadRequest)
		return
	}

	dir, err := NewWorkDir("sign-*")
	if err != nil {
		WriteError(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	in, err := saveUpload(files[0], dir, "sign", -1, []uploadedKind{pdfKind}, "PDF")
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}
	img, err := saveUpload(marks[0], dir, "mark", -1, []uploadedKind{pngKind}, "gambar PNG")
	if err != nil {
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := filepath.Join(dir, "signed.pdf")
	processErr := pdf.Sign(r.Context(), in, out, img,
		r.FormValue("page"), r.FormValue("x"), r.FormValue("y"), r.FormValue("width"))
	if processErr != nil {
		if errors.Is(processErr, pdf.ErrParam) {
			WriteError(w, processErr.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("penandatanganan gagal", "err", processErr)
		WriteError(w, "tanda tangan tidak bisa ditambahkan", http.StatusUnprocessableEntity)
		return
	}
	if err := WriteFileResponse(w, out, "signed.pdf", "application/pdf"); err != nil {
		WriteError(w, "gagal mengirim hasil", http.StatusInternalServerError)
	}
}