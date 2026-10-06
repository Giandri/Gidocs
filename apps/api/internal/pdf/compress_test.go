package pdf

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// hasGhostscript memastikan bin Ghostscript tersedia untuk tes integrasi.
func hasGhostscript(t *testing.T) {
	t.Helper()
	bin := os.Getenv("GS_BIN")
	if bin == "" {
		bin = "gswin64c"
		if _, err := exec.LookPath(bin); err != nil {
			bin = "gs"
		}
	}
	if _, err := exec.LookPath(bin); err != nil {
		t.Skip("Ghostscript tidak tersedia, tes dilewati")
	}
}

func TestCompressWhitelist(t *testing.T) {
	ctx := context.Background()
	for _, level := range []string{"", "HIGH", "sedang", "1"} {
		err := Compress(ctx, "in.pdf", "out.pdf", level)
		if !errors.Is(err, ErrParam) {
			t.Errorf("level %q: err = %v, ingin ErrParam", level, err)
		}
	}
}

func TestCompress(t *testing.T) {
	hasGhostscript(t)
	ctx := context.Background()
	in := filepath.Join("..", "..", "testdata", "a.pdf")
	out := filepath.Join(t.TempDir(), "out.pdf")

	if err := Compress(ctx, in, out, "medium"); err != nil {
		t.Fatalf("Compress: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("baca hasil: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("hasil bukan PDF: %.8q", data[:min(8, len(data))])
	}
}

func TestCompressMissingBinary(t *testing.T) {
	if os.Getenv("GS_BIN") == "" {
		t.Setenv("GS_BIN", "ghostscript-tidak-ada-xyz")
	} else {
		t.Skip("GS_BIN di-set dari luar")
	}
	ctx := context.Background()
	err := Compress(ctx, filepath.Join("..", "..", "testdata", "a.pdf"),
		filepath.Join(t.TempDir(), "out.pdf"), "medium")
	if err == nil || errors.Is(err, ErrParam) {
		t.Fatalf("ingin error generik karena bin tidak ada, dapat %v", err)
	}
}
