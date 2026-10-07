package pdf

import (
	"context"
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// copyFixture menyalin berkas dari testdata ke path tujuan.
func copyFixture(t *testing.T, dst, name string) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", name))
	if err != nil {
		t.Skipf("testdata hilangan: %v", err)
	}
	if err := os.WriteFile(dst, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// writePNG menulis PNG berukuran w x h piksel untuk dipakai sebagai tanda tangan.
func writePNG(t *testing.T, path string, w, h int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	img.Set(w/2, h/2, color.RGBA{R: 0, A: 255})
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

// writeInkBandPNG menulis PNG dengan tinta hanya di pita atas, berguna untuk
// memeriksa orientasi gambar setelah ditempel.
func writeInkBandPNG(t *testing.T, path string, w, h int) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h/4; y++ {
			img.Set(x, y, color.RGBA{A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestSignParams(t *testing.T) {
	bad := []struct {
		name              string
		page, x, y, width string
	}{
		{"halaman kosong", "", "100", "100", "150"},
		{"halaman bukan angka", "abc", "100", "100", "150"},
		{"halaman nol", "0", "100", "100", "150"},
		{"x bukan angka", "1", "abc", "100", "150"},
		{"x negatif", "1", "-5", "100", "150"},
		{"x di luar halaman", "1", "5000", "100", "150"},
		{"y negatif", "1", "100", "-5", "150"},
		{"y di luar halaman", "1", "100", "9000", "150"},
		{"lebar kosong", "1", "100", "100", ""},
		{"lebar nol", "1", "100", "100", "0"},
		{"lebar terlalu kecil", "1", "100", "100", "5"},
		{"lebar terlalu besar", "1", "100", "100", "5000"},
	}
	for _, c := range bad {
		t.Run(c.name, func(t *testing.T) {
			// PNG valid supaya kegagalan murni dari parameter, bukan berkas.
			img := filepath.Join(t.TempDir(), "sig.png")
			writePNG(t, img, 300, 100)
			if _, err := signDesc(img, c.page, c.x, c.y, c.width); !errors.Is(err, ErrParam) {
				t.Errorf("signDesc(%q,%q,%q,%q) = %v, ingin ErrParam", c.page, c.x, c.y, c.width, err)
			}
		})
	}
}

func TestSignDescAbsoluteTopLeft(t *testing.T) {
	img := filepath.Join(t.TempDir(), "sig.png")
	writePNG(t, img, 300, 100)

	// Gambar 300x100 px, lebar target 150 pt -> skala absolut 0.5.
	desc, err := signDesc(img, "2", "72", "144", "150")
	if err != nil {
		t.Fatalf("signDesc: %v", err)
	}
	if !strings.Contains(desc, "position:tl") {
		t.Errorf("desc %q harus memakai anchor top-left", desc)
	}
	if !strings.Contains(desc, "offset:72 -144") {
		t.Errorf("desc %q harus menggeser 72pt dari kiri dan 144pt dari atas", desc)
	}
	if !strings.Contains(desc, "scalefactor:0.50 abs") {
		t.Errorf("desc %q harus memakai skala absolut", desc)
	}
	// pdfcpu memutar watermark sebesar sudut diagonal halaman kecuali rotasi
	// dinyatakan eksplisit, jadi "rotation:0" wajib ada.
	if !strings.Contains(desc, "rotation:0") {
		t.Errorf("desc %q harus berisi rotation:0 agar gambar tidak miring", desc)
	}
}

func TestSignImageWidthFromFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sig.png")
	writePNG(t, path, 300, 100)

	width, err := imageWidthPt(path)
	if err != nil {
		t.Fatalf("imageWidthPt: %v", err)
	}
	if width != 300 {
		t.Errorf("lebar piksel = %d, ingin 300", width)
	}
}

func TestSignImageWidthRejectsNonPNG(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sig.png")
	if err := os.WriteFile(path, []byte("bukan png"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := imageWidthPt(path); !errors.Is(err, ErrParam) {
		t.Errorf("err = %v, ingin ErrParam untuk file bukan PNG", err)
	}
}

func TestSignWritesSignedPDF(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.pdf")
	copyFixture(t, in, "multi.pdf")
	img := filepath.Join(dir, "sig.png")
	writePNG(t, img, 300, 100)
	out := filepath.Join(dir, "out.pdf")

	if err := Sign(context.Background(), in, out, img, "2", "72", "144", "150"); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("hasil tidak ada: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("hasil bukan PDF: %.8q", data[:min(8, len(data))])
	}
	// Jumlah halaman harus tetap sama seperti aslinya (5 halaman).
	count, err := api.PageCountFile(context.Background(), out)
	if err != nil {
		t.Fatalf("PageCountFile: %v", err)
	}
	if count != 5 {
		t.Errorf("jumlah halaman = %d, ingin 5", count)
	}
}

func TestSignPageOutOfRange(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.pdf")
	copyFixture(t, in, "multi.pdf")
	img := filepath.Join(dir, "sig.png")
	writePNG(t, img, 300, 100)

	err := Sign(context.Background(), in, filepath.Join(dir, "out.pdf"), img, "9", "72", "144", "150")
	if !errors.Is(err, ErrParam) {
		t.Errorf("err = %v, ingin ErrParam untuk halaman di luar jangkauan", err)
	}
}