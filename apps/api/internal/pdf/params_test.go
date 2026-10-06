package pdf

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// TestParamWhitelist memastikan semua parameter masuk dicek whitelist
// dan mengembalikan ErrParam tanpa menyentuh file sama sekali.
func TestParamWhitelist(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		run  func() error
	}{
		{"rotate angle", func() error { return Rotate(ctx, "x.pdf", "y.pdf", "45", "") }},
		{"rotate angle teks", func() error { return Rotate(ctx, "x.pdf", "y.pdf", "abc", "") }},
		{"rotate pages teks", func() error { return Rotate(ctx, "x.pdf", "y.pdf", "90", "halaman1") }},
		{"rotate pages karakter", func() error { return Rotate(ctx, "x.pdf", "y.pdf", "90", "1;x") }},
		{"remove tanpa halaman", func() error { return RemovePages(ctx, "x.pdf", "y.pdf", "") }},
		{"remove halaman buruk", func() error { return RemovePages(ctx, "x.pdf", "y.pdf", "1;x") }},
		{"reorder kosong", func() error { return Reorder(ctx, "x.pdf", "y.pdf", "") }},
		{"reorder teks", func() error { return Reorder(ctx, "x.pdf", "y.pdf", "halaman") }},
		{"split mode asing", func() error { return Split(ctx, "x.pdf", ".", "chunks", "") }},
		{"split ranges kosong", func() error { return Split(ctx, "x.pdf", ".", "ranges", "") }},
		{"split ranges buruk", func() error { return Split(ctx, "x.pdf", ".", "ranges", "1-3,a") }},
		{"protect tanpa password", func() error { return Protect(ctx, "x.pdf", "y.pdf", "") }},
		{"unlock password panjang", func() error {
			return Unlock(ctx, "x.pdf", "y.pdf", strings.Repeat("a", 129))
		}},
		{"watermark teks kosong", func() error { return Watermark(ctx, "x.pdf", "y.pdf", " ", "center", "50") }},
		{"watermark posisi asing", func() error { return Watermark(ctx, "x.pdf", "y.pdf", "hi", "middle", "50") }},
		{"watermark opacity rendah", func() error { return Watermark(ctx, "x.pdf", "y.pdf", "hi", "center", "5") }},
		{"watermark opacity tinggi", func() error { return Watermark(ctx, "x.pdf", "y.pdf", "hi", "center", "500") }},
		{"watermark opacity teks", func() error {
			return Watermark(ctx, "x.pdf", "y.pdf", "hi", "center", "lima")
		}},
		{"images ukuran asing", func() error { return ImagesToPDF(ctx, []string{"a.png"}, "y.pdf", "letter") }},
		{"images kosong", func() error { return ImagesToPDF(ctx, nil, "y.pdf", "a4") }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := c.run(); !errors.Is(err, ErrParam) {
				t.Fatalf("expected ErrParam, got %v", err)
			}
		})
	}
}

func TestPageSelectionOpenRange(t *testing.T) {
	sel, err := pageSelection("1-,3,5-7")
	if err != nil {
		t.Fatalf("rentang terbuka harus diterima: %v", err)
	}
	if len(sel) != 3 {
		t.Fatalf("expected 3 bagian, got %d: %v", len(sel), sel)
	}
}
