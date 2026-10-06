package handler

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Giandri/Gidocs/apps/api/internal/jobs"
	"github.com/go-chi/chi/v5"
)

// redisAvailable memeriksa Redis bisa dihubungi untuk tes integrasi.
func redisAvailable() bool {
	conn, err := net.DialTimeout("tcp", "localhost:6379", 300*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// jobRouter menyalin route jobs dari main untuk tes.
func jobRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/api/jobs", EnqueueJob)
	r.Get("/api/jobs/{id}", JobStatus)
	r.Get("/api/jobs/{id}/download", JobDownload)
	return r
}

func getJob(t *testing.T, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	jobRouter().ServeHTTP(rec, req)
	return rec
}

func TestEnqueueJobValidation(t *testing.T) {
	doc := upload{name: "a.docx", content: string([]byte{'P', 'K', 0x03, 0x04, 0x14, 0x00, 0x06, 0x00})}

	// tanpa file
	assertCode(t, post(EnqueueJob, nil, map[string]string{"tool": "docx-to-pdf"}), http.StatusBadRequest)
	// tool tak dikenal
	assertCode(t, post(EnqueueJob, []upload{doc}, map[string]string{"tool": "merge-pdf"}), http.StatusBadRequest)
	// dua file
	assertCode(t, post(EnqueueJob, []upload{doc, doc}, map[string]string{"tool": "docx-to-pdf"}), http.StatusBadRequest)
	// magic bytes bukan ZIP (PDF)
	assertCode(t, post(EnqueueJob, []upload{fixture(t, "a.pdf")}, map[string]string{"tool": "docx-to-pdf"}), http.StatusBadRequest)
}

func TestEnqueueJobSuccess(t *testing.T) {
	if !redisAvailable() {
		t.Skip("Redis tidak tersedia, tes dilewati")
	}
	doc := upload{name: "a.docx", content: string([]byte{'P', 'K', 0x03, 0x04, 0x14, 0x00, 0x06, 0x00})}
	rec := post(EnqueueJob, []upload{doc}, map[string]string{"tool": "docx-to-pdf"})
	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !jobs.ValidID(resp.ID) || resp.Status != "queued" {
		t.Fatalf("respons aneh: %+v", resp)
	}
	defer os.RemoveAll(jobs.JobDir(resp.ID))

	meta, err := jobs.GetMeta(t.Context(), resp.ID)
	if err != nil {
		t.Fatalf("meta belum tersimpan: %v", err)
	}
	if meta.Status != "queued" || meta.Tool != "docx-to-pdf" {
		t.Fatalf("meta aneh: %+v", meta)
	}

	// download sebelum selesai = 409
	assertCode(t, getJob(t, "/api/jobs/"+resp.ID+"/download"), http.StatusConflict)
}

func TestJobStatusInvalidID(t *testing.T) {
	assertCode(t, getJob(t, "/api/jobs/bukan-id-valid"), http.StatusBadRequest)
	if redisAvailable() {
		assertCode(t, getJob(t, "/api/jobs/00000000000000000000000000000000"), http.StatusNotFound)
	}
}
