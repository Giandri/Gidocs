package jobs

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// copyTestPDF menyalin PDF sampel testdata ke folder job lalu mengembalikan pathnya.
func copyTestPDF(t *testing.T, dir string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "multi.pdf"))
	if err != nil {
		t.Skipf("testdata hilangan: %v", err)
	}
	dst := filepath.Join(dir, "input.pdf")
	if err := os.WriteFile(dst, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return dst
}

func TestConvertToImages(t *testing.T) {
	if _, err := FindPDFTOPPM(); err != nil {
		t.Skip("Poppler tidak tersedia, tes dilewati")
	}
	id := NewID()
	dir := JobDir(id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	input := copyTestPDF(t, dir)

	payload := Payload{JobID: id, Tool: "pdf-to-image", Input: input, OutDir: dir}
	if err := ConvertToImages(context.Background(), payload); err != nil {
		t.Fatalf("ConvertToImages: %v", err)
	}

	zr, err := zip.OpenReader(ZipResultPath(id))
	if err != nil {
		t.Fatalf("hasil bukan zip yang valid: %v", err)
	}
	defer zr.Close()
	pngs := 0
	for _, f := range zr.File {
		if filepath.Ext(f.Name) != ".png" {
			continue
		}
		pngs++
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("membuka %s: %v", f.Name, err)
		}
		n, _ := io.Copy(io.Discard, rc)
		rc.Close()
		if n == 0 {
			t.Errorf("%s kosong", f.Name)
		}
	}
	if pngs != 5 {
		t.Fatalf("jumlah png = %d, ingin 5 (multi.pdf 5 halaman)", pngs)
	}
}
