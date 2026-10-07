package handler

import (
	"archive/zip"
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

type upload struct{ name, content string }

func post(h http.HandlerFunc, files []upload, fields map[string]string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	for _, f := range files {
		part, _ := mw.CreateFormFile("files", f.name)
		part.Write([]byte(f.content))
	}
	for k, v := range fields {
		mw.WriteField(k, v)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/test", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

// postMark seperti post, tapi menambahkan satu file tambahan di field "mark".
func postMark(h http.HandlerFunc, files []upload, mark upload, fields map[string]string) *httptest.ResponseRecorder {
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	for _, f := range files {
		part, _ := mw.CreateFormFile("files", f.name)
		part.Write([]byte(f.content))
	}
	if mark.name != "" {
		part, _ := mw.CreateFormFile("mark", mark.name)
		part.Write([]byte(mark.content))
	}
	for k, v := range fields {
		mw.WriteField(k, v)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/test", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func fixture(t *testing.T, name string) upload {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "testdata", name))
	if err != nil {
		t.Skipf("testdata hilangan: %v", err)
	}
	return upload{name: name, content: string(b)}
}

func imageUpload(t *testing.T, kind string) upload {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	var err error
	if kind == "png" {
		err = png.Encode(&buf, img)
		if err == nil {
			return upload{name: "img.png", content: buf.String()}
		}
	} else {
		err = jpeg.Encode(&buf, img, nil)
		if err == nil {
			return upload{name: "img.jpg", content: buf.String()}
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	return upload{}
}

func assertPDF(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte("%PDF-")) {
		t.Fatal("respons bukan PDF")
	}
}

func assertCode(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("expected %d, got %d: %s", want, rec.Code, rec.Body.String())
	}
}

func pageCount(t *testing.T, rec *httptest.ResponseRecorder) int {
	t.Helper()
	path := filepath.Join(t.TempDir(), "out.pdf")
	if err := os.WriteFile(path, rec.Body.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
	n, err := api.PageCountFile(context.Background(), path)
	if err != nil {
		t.Fatalf("gagal membaca hasil: %v", err)
	}
	return n
}

func TestRotateHandler(t *testing.T) {
	multi := fixture(t, "multi.pdf")

	rec := post(RotateHandler, []upload{multi}, map[string]string{"angle": "90"})
	assertPDF(t, rec)

	assertCode(t, post(RotateHandler, []upload{multi}, map[string]string{"angle": "45"}), http.StatusBadRequest)
	assertCode(t, post(RotateHandler, []upload{multi}, map[string]string{"angle": ""}), http.StatusBadRequest)
	assertCode(t, post(RotateHandler, []upload{multi, multi}, map[string]string{"angle": "90"}), http.StatusBadRequest)
	assertCode(t, post(RotateHandler, []upload{{name: "x.txt", content: "bukan pdf"}}, map[string]string{"angle": "90"}), http.StatusBadRequest)
	assertCode(t, post(RotateHandler, nil, map[string]string{"angle": "90"}), http.StatusBadRequest)
}

func TestSplitHandlerEveryPage(t *testing.T) {
	rec := post(SplitHandler, []upload{fixture(t, "multi.pdf")}, map[string]string{"mode": "every"})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/zip" {
		t.Fatalf("expected application/zip, got %s", ct)
	}

	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("respons bukan zip: %v", err)
	}
	if len(zr.File) != 5 {
		t.Fatalf("expected 5 bagian, got %d", len(zr.File))
	}
}

func TestSplitHandlerRanges(t *testing.T) {
	rec := post(SplitHandler, []upload{fixture(t, "multi.pdf")}, map[string]string{
		"mode": "ranges", "ranges": "1-2,4",
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	zr, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("respons bukan zip: %v", err)
	}
	if len(zr.File) != 2 {
		t.Fatalf("expected 2 bagian, got %d", len(zr.File))
	}
}

func TestSplitHandlerInvalidRanges(t *testing.T) {
	assertCode(t, post(SplitHandler, []upload{fixture(t, "multi.pdf")}, map[string]string{
		"mode": "ranges", "ranges": "abc",
	}), http.StatusBadRequest)

	assertCode(t, post(SplitHandler, []upload{fixture(t, "multi.pdf")}, map[string]string{
		"mode": "ranges", "ranges": "99-100",
	}), http.StatusBadRequest)
}

func TestRemovePagesHandler(t *testing.T) {
	multi := fixture(t, "multi.pdf")

	rec := post(RemovePagesHandler, []upload{multi}, map[string]string{"pages": "1,3"})
	assertPDF(t, rec)
	if got := pageCount(t, rec); got != 3 {
		t.Fatalf("expected 3 halaman, got %d", got)
	}

	assertCode(t, post(RemovePagesHandler, []upload{multi}, map[string]string{"pages": ""}), http.StatusBadRequest)
	assertCode(t, post(RemovePagesHandler, []upload{multi}, map[string]string{"pages": "1;x"}), http.StatusBadRequest)
}

func TestReorderHandler(t *testing.T) {
	multi := fixture(t, "multi.pdf")

	rec := post(ReorderHandler, []upload{multi}, map[string]string{"order": "5,1,2"})
	assertPDF(t, rec)
	if got := pageCount(t, rec); got != 3 {
		t.Fatalf("expected 3 halaman, got %d", got)
	}

	assertCode(t, post(ReorderHandler, []upload{multi}, map[string]string{"order": ""}), http.StatusBadRequest)
}

func TestProtectThenUnlock(t *testing.T) {
	a := fixture(t, "a.pdf")

	rec := post(ProtectHandler, []upload{a}, map[string]string{"password": "rahasia123"})
	assertPDF(t, rec)
	enc := upload{name: "enc.pdf", content: rec.Body.String()}

	assertPDF(t, post(UnlockHandler, []upload{enc}, map[string]string{"password": "rahasia123"}))

	// password salah
	assertCode(t, post(UnlockHandler, []upload{enc}, map[string]string{"password": "salah"}), http.StatusUnprocessableEntity)
	// file tidak terproteksi
	assertCode(t, post(UnlockHandler, []upload{a}, map[string]string{"password": "x"}), http.StatusUnprocessableEntity)
	// protect tanpa password
	assertCode(t, post(ProtectHandler, []upload{a}, map[string]string{"password": ""}), http.StatusBadRequest)
}

func TestWatermarkHandler(t *testing.T) {
	a := fixture(t, "a.pdf")

	rec := post(WatermarkHandler, []upload{a}, map[string]string{
		"text": "DRAFT", "position": "diagonal", "opacity": "50",
	})
	assertPDF(t, rec)

	assertCode(t, post(WatermarkHandler, []upload{a}, map[string]string{
		"text": "", "position": "center", "opacity": "50",
	}), http.StatusBadRequest)
	assertCode(t, post(WatermarkHandler, []upload{a}, map[string]string{
		"text": "DRAFT", "position": "atas", "opacity": "50",
	}), http.StatusBadRequest)
assertCode(t, post(WatermarkHandler, []upload{a}, map[string]string{
		"text": "DRAFT", "position": "center", "opacity": "5",
	}), http.StatusBadRequest)
}

func TestWatermarkImageHandler(t *testing.T) {
	a := fixture(t, "a.pdf")
	pngUp := imageUpload(t, "png")

	rec := postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "center", "opacity": "50",
	})
	assertPDF(t, rec)

	// gambar bukan PNG/JPEG
	assertCode(t, postMark(WatermarkHandler, []upload{a}, upload{name: "x.txt", content: "bukan gambar"}, map[string]string{
		"position": "center", "opacity": "50",
	}), http.StatusBadRequest)

	// diagonal khusus teks, tidak valid untuk gambar
	assertCode(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "diagonal", "opacity": "50",
	}), http.StatusBadRequest)

	// posisi tidak dikenal
	assertCode(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "atas", "opacity": "50",
	}), http.StatusBadRequest)

	// opasitas di luar rentang
	assertCode(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "center", "opacity": "5",
	}), http.StatusBadRequest)

	// ukuran di luar rentang
	assertCode(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "center", "opacity": "50", "size": "5",
	}), http.StatusBadRequest)
	assertCode(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "center", "opacity": "50", "size": "105",
	}), http.StatusBadRequest)

	// ukuran berbeda wajar
	assertPDF(t, postMark(WatermarkHandler, []upload{a}, pngUp, map[string]string{
		"position": "center", "opacity": "50", "size": "80",
	}))
}

func TestImagesToPDFHandler(t *testing.T) {
	pngUp := imageUpload(t, "png")
	jpgUp := imageUpload(t, "jpg")

	rec := post(ImagesToPDFHandler, []upload{pngUp, jpgUp}, map[string]string{"size": "a4"})
	assertPDF(t, rec)
	if got := pageCount(t, rec); got != 2 {
		t.Fatalf("expected 2 halaman, got %d", got)
	}

	rec = post(ImagesToPDFHandler, []upload{pngUp}, map[string]string{"size": "fit"})
	assertPDF(t, rec)

	assertCode(t, post(ImagesToPDFHandler, []upload{{name: "x.txt", content: "bukan gambar"}}, map[string]string{"size": "a4"}), http.StatusBadRequest)
	assertCode(t, post(ImagesToPDFHandler, []upload{pngUp}, map[string]string{"size": "letter"}), http.StatusBadRequest)
	assertCode(t, post(ImagesToPDFHandler, nil, map[string]string{"size": "a4"}), http.StatusBadRequest)
}

func ghostscriptAvailable() bool {
	bin := os.Getenv("GS_BIN")
	if bin == "" {
		bin = "gswin64c"
		if _, err := exec.LookPath(bin); err != nil {
			bin = "gs"
		}
	}
	_, err := exec.LookPath(bin)
	return err == nil
}

func TestCompressHandler(t *testing.T) {
	doc := fixture(t, "a.pdf")

	// Level di luar whitelist ditolak bahkan sebelum Ghostscript dijalankan.
	assertCode(t, post(CompressHandler, []upload{doc}, map[string]string{"level": "extreme"}), http.StatusBadRequest)
	assertCode(t, post(CompressHandler, []upload{doc}, nil), http.StatusBadRequest)

	if !ghostscriptAvailable() {
		t.Log("Ghostscript tidak tersedia, tes sukses dilewati")
		return
	}
	assertPDF(t, post(CompressHandler, []upload{doc}, map[string]string{"level": "medium"}))
}
