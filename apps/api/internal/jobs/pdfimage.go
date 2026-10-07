package jobs

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// FindPDFTOPPM mencari biner pdftoppm (Poppler).
// Urutan: env PDFTOPPM_BIN -> PATH -> lokasi umum scoop.
func FindPDFTOPPM() (string, error) {
	if v := os.Getenv("PDFTOPPM_BIN"); v != "" {
		return v, nil
	}
	if p, err := exec.LookPath("pdftoppm"); err == nil {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err == nil {
		if p, _ := filepath.Glob(filepath.Join(home, "scoop", "apps", "poppler", "current", "bin", "pdftoppm.exe")); len(p) > 0 {
			return p[0], nil
		}
	}
	return "", errors.New("Poppler (pdftoppm) tidak ditemukan")
}

// ConvertToImages merender tiap halaman PDF menjadi PNG (150 dpi) lalu
// mengemasnya ke ZipResultPath(JobID).
func ConvertToImages(ctx context.Context, p Payload) error {
	bin, err := FindPDFTOPPM()
	if err != nil {
		return err
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	prefix := filepath.Join(p.OutDir, "page")
	cmd := exec.CommandContext(runCtx, bin, "-png", "-r", "150", p.Input, prefix)
	var output bytes.Buffer
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		// output hanya untuk log internal, tidak pernah dikirim ke user.
		slog.Error("pdftoppm gagal", "err", err, "output", output.String())
		return errors.New("konversi PDF ke gambar gagal, periksa file PDF")
	}

	pages, err := filepath.Glob(prefix + "-*.png")
	if err != nil || len(pages) == 0 {
		return errors.New("Poppler tidak menghasilkan gambar")
	}
	if err := zipFiles(pages, ZipResultPath(p.JobID)); err != nil {
		return err
	}
	return nil
}

// zipFiles mengemas daftar file menjadi satu arsip ZIP di dest.
// Hanya dipanggil dengan path buatan server (bukan input user).
func zipFiles(paths []string, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("gagal membuat zip: %w", err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, path := range paths {
		w, err := zw.Create(filepath.Base(path))
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
	}
	return zw.Close()
}
