package jobs

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// officeTools memetakan slug tool ke ekstensi input yang diterima (whitelist).
var officeTools = map[string]string{
	"docx-to-pdf":  ".docx",
	"excel-to-pdf": ".xlsx",
	"pptx-to-pdf":  ".pptx",
}

// OfficeExt mengembalikan ekstensi (pakai titik) untuk slug Office.
func OfficeExt(slug string) (string, bool) {
	ext, ok := officeTools[slug]
	return ext, ok
}

// FindSOFFICE mencari biner LibreOffice.
// Urutan: env SOFFICE_BIN -> PATH (soffice.com dulu, karena di Windows
// hanya .com yang menunggu proses konversi selesai) -> lokasi umum.
func FindSOFFICE() (string, error) {
	if v := os.Getenv("SOFFICE_BIN"); v != "" {
		return v, nil
	}
	for _, name := range []string{"soffice.com", "soffice"} {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	home, err := os.UserHomeDir()
	if err == nil {
		if p, _ := filepath.Glob(filepath.Join(home, "scoop", "apps", "libreoffice*",
			"current", "LibreOffice", "program", "soffice.com")); len(p) > 0 {
			return p[0], nil
		}
	}
	if p, _ := filepath.Glob(`C:\Program Files\LibreOffice*\program\soffice.com`); len(p) > 0 {
		return p[0], nil
	}
	return "", errors.New("LibreOffice tidak ditemukan")
}

// ConvertOffice menjalankan LibreOffice headless untuk satu payload job.
// Hasil akhirnya adalah ResultPath(JobID) berupa PDF.
func ConvertOffice(ctx context.Context, p Payload) error {
	bin, err := FindSOFFICE()
	if err != nil {
		return err
	}
	runCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Profil unik per job agar beberapa konversi tidak saling mengunci.
	profileURL := "file:///" + filepath.ToSlash(filepath.Join(p.OutDir, "lo-profile"))

	cmd := exec.CommandContext(runCtx, bin,
		"--headless", "--norestore", "--nolockcheck",
		"-env:UserInstallation="+profileURL,
		"--convert-to", "pdf",
		"--outdir", p.OutDir,
		p.Input,
	)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		// output hanya untuk log internal, tidak pernah dikirim ke user.
		slog.Error("libreoffice gagal", "err", err, "output", output.String())
		return errors.New("konversi gagal, format file tidak didukung")
	}

	generated := strings.TrimSuffix(p.Input, filepath.Ext(p.Input)) + ".pdf"
	if _, err := os.Stat(generated); err != nil {
		return errors.New("LibreOffice tidak menghasilkan PDF")
	}
	if err := os.Rename(generated, ResultPath(p.JobID)); err != nil {
		return fmt.Errorf("gagal menyimpan hasil: %w", err)
	}
	return nil
}
