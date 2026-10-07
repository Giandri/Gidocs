package handler

import (
	"net/http"
	"testing"
)

// TestSignHandlerValidation menguji penolakan input tidak valid untuk Sign PDF.
func TestSignHandlerValidation(t *testing.T) {
	pdf := fixture(t, "multi.pdf")
	png := imageUpload(t, "png")
	pdfFields := map[string]string{"page": "1", "x": "72", "y": "144", "width": "150"}

	// tanpa PDF
	assertCode(t, postMark(SignHandler, []upload{pdf}, png, map[string]string{}), http.StatusBadRequest)
	// tanpa tanda tangan
	assertCode(t, postMark(SignHandler, []upload{pdf}, upload{}, pdfFields), http.StatusBadRequest)
	// dua PDF
	assertCode(t, postMark(SignHandler, []upload{pdf, pdf}, png, pdfFields), http.StatusBadRequest)
	// PDF bukan PDF
	assertCode(t, postMark(SignHandler, []upload{{name: "x.txt", content: "bukan pdf"}}, png, pdfFields), http.StatusBadRequest)
	// tanda tangan bukan gambar
	assertCode(t, postMark(SignHandler, []upload{pdf}, upload{name: "sig.txt", content: "bukan gambar"}, pdfFields), http.StatusBadRequest)
	// parameter di luar jangkauan
	for _, fields := range []map[string]string{
		{"page": "0", "x": "72", "y": "144", "width": "150"},
		{"page": "99", "x": "72", "y": "144", "width": "150"},
		{"page": "1", "x": "-1", "y": "144", "width": "150"},
		{"page": "1", "x": "72", "y": "144", "width": "0"},
		{"page": "abc", "x": "72", "y": "144", "width": "150"},
	} {
		assertCode(t, postMark(SignHandler, []upload{pdf}, png, fields), http.StatusBadRequest)
	}
}

// TestSignHandlerSuccess memverifikasi PDF bertanda tangan bisa diunduh.
func TestSignHandlerSuccess(t *testing.T) {
	pdf := fixture(t, "multi.pdf")
	png := imageUpload(t, "png")
	rec := postMark(SignHandler, []upload{pdf}, png, map[string]string{
		"page": "1", "x": "72", "y": "144", "width": "150",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("Content-Type = %q, ingin application/pdf", ct)
	}
	body := rec.Body.Bytes()
	if len(body) < 5 || string(body[:5]) != "%PDF-" {
		t.Fatalf("body bukan PDF: %.8q", body[:min(8, len(body))])
	}
}