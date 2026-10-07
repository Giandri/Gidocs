package pdf

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Batas koordinat tanda tangan dalam poin PDF. Lebar dibatasi agar tanda
// tangan tidak menutupi seluruh halaman.
const (
	maxSignOffset = 2000 // pt dari tepi kiri/atas
	minSignWidth  = 40   // pt
	maxSignWidth  = 400  // pt
)

// imageWidthPt membaca lebar gambar tanda tangan dalam piksel. pdfcpu memakai
// 1 piksel = 1 poin untuk watermark gambar, jadi piksel bisa langsung dipakai
// untuk menghitung skala absolut.
func imageWidthPt(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("%w: tanda tangan tidak terbaca", ErrParam)
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil || cfg.Width <= 0 {
		return 0, fmt.Errorf("%w: tanda tangan bukan gambar yang valid", ErrParam)
	}
	return cfg.Width, nil
}

// signDesc memvalidasi parameter dan membangun deskripsi penempatan tanda tangan
// untuk pdfcpu: anchor top-left, offset dari tepi kiri/atas dalam pt, dan skala
// absolut supaya lebar tanda tangan pas dengan yang diminta.
//
// y diukur dari tepi atas halaman (ikut koordinat browser), sedangkan sumbu y
// pdfcpu ke atas, jadi y dibalik tandanya.
func signDesc(imgPath, page, x, y, width string) (string, error) {
	p, err := strconv.Atoi(page)
	if err != nil || p < 1 {
		return "", fmt.Errorf("%w: halaman harus angka mulai 1", ErrParam)
	}
	px, err := parseCoordinate(x, "posisi horizontal")
	if err != nil {
		return "", err
	}
	py, err := parseCoordinate(y, "posisi vertikal")
	if err != nil {
		return "", err
	}
	w, err := strconv.Atoi(width)
	if err != nil || w < minSignWidth || w > maxSignWidth {
		return "", fmt.Errorf("%w: lebar tanda tangan harus %d-%d pt", ErrParam, minSignWidth, maxSignWidth)
	}
	imgWidth, err := imageWidthPt(imgPath)
	if err != nil {
		return "", err
	}
	scale := float64(w) / float64(imgWidth)
	// rotation:0 wajib ditulis: tanpa itu pdfcpu memakai diagonal default dan
	// memutar gambar sebesar sudut halaman. Menyetel rotation juga mematikan diagonal.
	return fmt.Sprintf("position:tl,rotation:0,offset:%d %d,scalefactor:%s abs",
		px, -py, strconv.FormatFloat(scale, 'f', 2, 64)), nil
}

// parseCoordinate memvalidasi jarak dari tepi halaman dalam poin bulat.
func parseCoordinate(v, label string) (int, error) {
	n, err := strconv.Atoi(v)
	if err != nil || n < 0 || n > maxSignOffset {
		return 0, fmt.Errorf("%w: %s harus 0-%d pt", ErrParam, label, maxSignOffset)
	}
	return n, nil
}

// Sign menempelkan gambar tanda tangan ke satu halaman PDF pada koordinat
// tertentu (poin PDF, 1-based untuk page). Tinggi gambar mengikuti proporsi
// aslinya.
func Sign(ctx context.Context, in, out, imgPath, page, x, y, width string) error {
	sel, err := pageSelection(page)
	if err != nil {
		return err
	}
	if err := checkSelectionBounds(ctx, in, sel); err != nil {
		return err
	}
	desc, err := signDesc(imgPath, page, x, y, width)
	if err != nil {
		return err
	}
	return api.AddImageWatermarksFile(ctx, in, out, sel, true, imgPath, desc, nil)
}