package pdf

import (
	"context"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
)

// renderPage merender satu halaman PDF menjadi gambar grayscale lewat pdftoppm.
// Skala 2 piksel per poin.
func renderPage(t *testing.T, pdfPath string, page int) image.Image {
	t.Helper()
	bin, err := exec.LookPath("pdftoppm")
	if err != nil {
		t.Skip("Poppler tidak tersedia, tes dilewati")
	}
	prefix := filepath.Join(t.TempDir(), "page")
	// -gray untuk hitung piksel gelap, -f/-l untuk satu halaman, -r 144 = 2 px/pt.
	cmd := exec.Command(bin, "-png", "-gray", "-r", "144", "-f", strconv.Itoa(page),
		"-l", strconv.Itoa(page), pdfPath, prefix)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("pdftoppm gagal: %v: %s", err, output)
	}
	pagePNG := prefix + "-" + strconv.Itoa(page) + ".png"
	f, err := os.Open(pagePNG)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	return img
}

// darkPixels menghitung piksel yang jauh lebih gelap dari putih di dalam kotak.
func darkPixels(img image.Image, x, y, w, h int) int {
	bounds := img.Bounds()
	count := 0
	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if px < bounds.Min.X || px >= bounds.Max.X || py < bounds.Min.Y || py >= bounds.Max.Y {
				continue
			}
			if r, g, b, _ := img.At(px, py).RGBA(); (r+g+b)/3 < 0x8000 {
			count++
		}
		}
	}
	return count
}

// TestSignKeepsOrientation memastikan gambar tanda tangan tidak diputar.
// pdfcpu memakai diagonal default sehingga gambar ikut miring sebesar sudut
// halaman; tanda tangan harus tetap tegak lurus dengan halaman.
func TestSignKeepsOrientation(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.pdf")
	copyFixture(t, in, "multi.pdf")
	img := filepath.Join(dir, "sig.png")
	writeInkBandPNG(t, img, 300, 100)
	out := filepath.Join(dir, "out.pdf")

	// Lebar 180 pt -> tinggi 60 pt (rasio 3:1), di 72pt dari kiri, 200pt dari atas.
	const xPt, yPt, widthPt, scale = 72, 200, 180, 2
	boxX, boxY := xPt*scale, yPt*scale
	boxW, boxH := widthPt*scale, 60*scale

	before := renderPage(t, in, 1)
	if got := darkPixels(before, boxX, boxY, boxW, boxH); got != 0 {
		t.Skipf("PDF asal sudah gelap di area tanda tangan (%d piksel)", got)
	}
	if err := Sign(context.Background(), in, out, img, "1", "72", "200", "180"); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	after := renderPage(t, out, 1)

	// Tinta hanya mengisi pita atas gambar: bagian atas area tanda tangan harus
	// ada tinta, bagian bawah harus kosong. Kalau diputar, tinta menyebar ke
	// seluruh area.
	top := darkPixels(after, boxX, boxY, boxW, boxH/2)
	bottom := darkPixels(after, boxX, boxY+boxH/2, boxW, boxH/2)
	if top == 0 {
		t.Errorf("bagian atas area tanda tangan kosong (%d piksel); tanda tangan tidak terlihat", top)
	}
	if bottom > 0 {
		t.Errorf("bagian bawah area tanda tangan berisi %d piksel; gambar kemungkinan diputar", bottom)
	}
}

// TestSignPlacesSignatureOnRequestedSpot memastikan tanda tangan muncul di
// halaman dan koordinat yang diminta, bukan sekadar PDF yang valid.
func TestSignPlacesSignatureOnRequestedSpot(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.pdf")
	copyFixture(t, in, "multi.pdf")
	img := filepath.Join(dir, "sig.png")
	writeInkBandPNG(t, img, 300, 100)
	out := filepath.Join(dir, "out.pdf")

	// 150x50 pt pada 72pt dari kiri, 144pt dari atas, halaman 1.
	// pdftoppm -r 144 = 2 piksel per poin.
	const (
		xPt, yPt, widthPt = 72, 144, 150
		scale             = 2
	)
	boxX, boxY := xPt*scale, yPt*scale
	boxW, boxH := widthPt*scale, 50*scale

	beforePage1 := darkPixels(renderPage(t, in, 1), boxX, boxY, boxW, boxH)
	beforePage2 := darkPixels(renderPage(t, in, 2), boxX, boxY, boxW, boxH)

	if err := Sign(context.Background(), in, out, img, "1", "72", "144", "150"); err != nil {
		t.Fatalf("Sign: %v", err)
	}

	afterPage1 := darkPixels(renderPage(t, out, 1), boxX, boxY, boxW, boxH)
	afterPage2 := darkPixels(renderPage(t, out, 2), boxX, boxY, boxW, boxH)

	if afterPage1 <= beforePage1 {
		t.Errorf("halaman 1: piksel gelap di area tanda tangan tidak bertambah (%d -> %d); tanda tangan tidak muncul di posisi diminta",
			beforePage1, afterPage1)
	}
	if afterPage2 != beforePage2 {
		t.Errorf("halaman 2 ikut berubah (%d -> %d); halaman selain 1 tidak boleh disentuh",
			beforePage2, afterPage2)
	}
}