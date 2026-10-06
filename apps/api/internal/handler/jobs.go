package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/Giandri/Gidocs/apps/api/internal/jobs"
	"github.com/go-chi/chi/v5"
)

var zipMagic = []byte{'P', 'K', 0x03, 0x04}

// writeJSON menulis respons JSON sederhana.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

// EnqueueJob menerima 1 file Office + field "tool", lalu menaruhnya di antrian.
// Merespons 202 dengan id job; klien memantau lewat GET /api/jobs/{id}.
func EnqueueJob(w http.ResponseWriter, r *http.Request) {
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

	slug := r.FormValue("tool")
	ext, ok := jobs.OfficeExt(slug)
	if !ok {
		WriteError(w, "tool tidak dikenal", http.StatusBadRequest)
		return
	}

	id := jobs.NewID()
	dir := jobs.JobDir(id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		WriteError(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}

	// Magic ZIP (docx/xlsx/pptx sama); ekstensi diambil dari whitelist tool,
	// bukan dari nama file user.
	kinds := []uploadedKind{{magic: zipMagic, ext: ext}}
	label := strings.ToUpper(strings.TrimPrefix(ext, "."))
	input, err := saveUpload(files[0], dir, "input", -1, kinds, label)
	if err != nil {
		os.RemoveAll(dir)
		WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	payload := jobs.Payload{JobID: id, Tool: slug, Input: input, OutDir: dir}
	if err := jobs.Enqueue(r.Context(), payload); err != nil {
		os.RemoveAll(dir)
		WriteError(w, "layanan antrian sedang tidak tersedia", http.StatusServiceUnavailable)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]string{"id": id, "status": "queued"})
}

// JobStatus mengembalikan status job dalam JSON untuk polling.
func JobStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !jobs.ValidID(id) {
		WriteError(w, "id job tidak valid", http.StatusBadRequest)
		return
	}
	meta, err := jobs.GetMeta(r.Context(), id)
	if errors.Is(err, jobs.ErrNotFound) {
		WriteError(w, "job tidak ditemukan atau sudah kedaluwarsa", http.StatusNotFound)
		return
	}
	if err != nil {
		WriteError(w, "layanan antrian sedang tidak tersedia", http.StatusServiceUnavailable)
		return
	}
	resp := map[string]string{"id": id, "status": meta.Status}
	if meta.Error != "" {
		resp["error"] = meta.Error
	}
	if meta.Status == "done" {
		resp["download"] = "/api/jobs/" + id + "/download"
	}
	writeJSON(w, http.StatusOK, resp)
}

// JobDownload mengirim PDF hasil konversi.
func JobDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if !jobs.ValidID(id) {
		WriteError(w, "id job tidak valid", http.StatusBadRequest)
		return
	}
	meta, err := jobs.GetMeta(r.Context(), id)
	if errors.Is(err, jobs.ErrNotFound) {
		WriteError(w, "job tidak ditemukan atau sudah kedaluwarsa", http.StatusNotFound)
		return
	}
	if err != nil {
		WriteError(w, "layanan antrian sedang tidak tersedia", http.StatusServiceUnavailable)
		return
	}
	if meta.Status == "error" {
		WriteError(w, meta.Error, http.StatusConflict)
		return
	}
	if meta.Status != "done" {
		WriteError(w, "job belum selesai", http.StatusConflict)
		return
	}
	if err := WriteFileResponse(w, jobs.ResultPath(id), "converted.pdf", "application/pdf"); err != nil {
		WriteError(w, "hasil sudah tidak tersedia", http.StatusNotFound)
	}
}
