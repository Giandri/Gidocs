package pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// gsPresets memetakan level UI (light/medium/strong) ke preset Ghostscript.
// Hanya nilai whitelist ini yang pernah masuk ke argumen perintah.
var gsPresets = map[string]string{
	"light":  "/printer", // jaga kualitas, kompresi ringan
	"medium": "/ebook",   // seimbang
	"strong": "/screen",  // agresif, ukuran terkecil
}

// Compress menjalankan Ghostscript dengan -dSAFER dan argumen tetap.
// level harus salah satu kunci gsPresets; semua path adalah milik server.
func Compress(ctx context.Context, in, out, level string) error {
	preset, ok := gsPresets[level]
	if !ok {
		return fmt.Errorf("%w: level harus light, medium, atau strong", ErrParam)
	}

	bin := os.Getenv("GS_BIN")
	if bin == "" {
		if runtime.GOOS == "windows" {
			bin = "gswin64c"
		} else {
			bin = "gs"
		}
	}

	runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Ghostscript (SAFER) hanya menerima path absolut.
	inAbs, err := filepath.Abs(in)
	if err != nil {
		return errors.New("path input tidak valid")
	}
	outAbs, err := filepath.Abs(out)
	if err != nil {
		return errors.New("path keluaran tidak valid")
	}

	cmd := exec.CommandContext(runCtx, bin,
		"-dSAFER", "-dBATCH", "-dNOPAUSE", "-dNOPROMPT",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS="+preset,
		"-o", outAbs,
		inAbs,
	)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		// output hanya untuk log internal, tidak pernah dikirim ke user.
		slog.Error("ghostscript gagal", "err", err, "output", output.String())
		return errors.New("ghostscript gagal memproses PDF")
	}
	if _, err := os.Stat(out); err != nil {
		return errors.New("ghostscript tidak menghasilkan file")
	}
	return nil
}
