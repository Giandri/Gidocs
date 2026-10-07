package jobs

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// FindOCRMYPDF mencari biner OCRmyPDF.
// Urutan: env OCRMYPDF_BIN -> PATH -> instalasi pip user.
func FindOCRMYPDF() (string, error) {
	if v := os.Getenv("OCRMYPDF_BIN"); v != "" {
		return v, nil
	}
	if p, err := exec.LookPath("ocrmypdf"); err == nil {
		return p, nil
	}
	if appdata := os.Getenv("APPDATA"); appdata != "" {
		if p, _ := filepath.Glob(filepath.Join(appdata, "Python", "*", "Scripts", "ocrmypdf.exe")); len(p) > 0 {
			return p[0], nil
		}
	}
	return "", errors.New("OCRmyPDF tidak ditemukan")
}

// findTessdata mengembalikan folder traineddata Tesseract bila TESSDATA_PREFIX
// belum diset (mis. terminal lama setelah pemasangan bahasa). Mengembalikan
// string kosong kalau env sudah ada atau tidak ditemukan.
func findTessdata() string {
	if os.Getenv("TESSDATA_PREFIX") != "" {
		return ""
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	p, _ := filepath.Glob(filepath.Join(home, "scoop", "apps", "tesseract-languages", "current"))
	if len(p) > 0 {
		return p[0]
	}
	return ""
}

// OCRPDF menjalankan OCRmyPDF dengan bahasa ind+eng dan menulis PDF hasil
// (teks tersembunyi) ke ResultPath(JobID).
func OCRPDF(ctx context.Context, p Payload) error {
	bin, err := FindOCRMYPDF()
	if err != nil {
		return err
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// --skip-text: halaman yang sudah ber-teks tidak di-OCR ulang.
	cmd := exec.CommandContext(runCtx, bin,
		"--skip-text", "--language", "ind+eng",
		p.Input, ResultPath(p.JobID))
	if dir := findTessdata(); dir != "" {
		cmd.Env = append(os.Environ(), "TESSDATA_PREFIX="+dir)
	}
	var output bytes.Buffer
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		// output hanya untuk log internal, tidak pernah dikirim ke user.
		slog.Error("ocrmypdf gagal", "err", err, "output", output.String())
		return errors.New("OCR gagal, periksa file PDF")
	}
	if _, err := os.Stat(ResultPath(p.JobID)); err != nil {
		return errors.New("OCRmyPDF tidak menghasilkan PDF")
	}
	return nil
}
