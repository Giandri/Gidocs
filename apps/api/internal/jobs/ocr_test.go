package jobs

import (
	"context"
	"os"
	"testing"
)

func TestOCRPDF(t *testing.T) {
	if _, err := FindOCRMYPDF(); err != nil {
		t.Skip("OCRmyPDF tidak tersedia, tes dilewati")
	}
	id := NewID()
	dir := JobDir(id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	input := copyTestPDF(t, dir)

	payload := Payload{JobID: id, Tool: "ocr-pdf", Input: input, OutDir: dir}
	if err := OCRPDF(context.Background(), payload); err != nil {
		t.Fatalf("OCRPDF: %v", err)
	}
	data, err := os.ReadFile(ResultPath(id))
	if err != nil {
		t.Fatalf("hasil tidak ada: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("hasil bukan PDF: %.8q", data[:min(8, len(data))])
	}
}
