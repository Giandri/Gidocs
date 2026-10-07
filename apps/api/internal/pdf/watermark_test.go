package pdf

import (
	"errors"
	"strings"
	"testing"
)

func TestSizeFactor(t *testing.T) {
	cases := []struct {
		in   string
		want float64
		err  bool
	}{
		{"", 0.5, false}, // kosong = default 50
		{"50", 0.5, false},
		{"10", 0.1, false},
		{"100", 1, false},
		{"5", 0, true},
		{"101", 0, true},
		{"abc", 0, true},
	}
	for _, c := range cases {
		got, err := sizeFactor(c.in)
		if c.err {
			if err == nil {
				t.Fatalf("sizeFactor(%q) harus error", c.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("sizeFactor(%q): %v", c.in, err)
		}
		if got != c.want {
			t.Fatalf("sizeFactor(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestWatermarkImageParams(t *testing.T) {
	valid := []string{"center", "top-left", "bottom-right"}
	for _, pos := range valid {
		if _, err := watermarkDesc(pos, "50", false); err != nil {
			t.Fatalf("posisi %s harus valid: %v", pos, err)
		}
	}

	// diagonal tidak boleh untuk watermark gambar
	if _, err := watermarkDesc("diagonal", "50", false); !errors.Is(err, ErrParam) {
		t.Fatal("diagonal untuk gambar harus ditolak")
	}
	// diagonal tetap valid untuk teks
	if _, err := watermarkDesc("diagonal", "50", true); err != nil {
		t.Fatalf("diagonal untuk teks harus valid: %v", err)
	}
	// opacity di luar rentang
	for _, o := range []string{"5", "101", "abc"} {
		if _, err := watermarkDesc("center", o, true); !errors.Is(err, ErrParam) {
			t.Fatalf("opacity %q harus ditolak", o)
		}
	}
	// posisi tidak dikenal
	if _, err := watermarkDesc("atas", "50", true); !errors.Is(err, ErrParam) {
		t.Fatal("posisi tidak dikenal harus ditolak")
	}
	// posisi + opacity menghasilkan desc berawalan posisi dan berakhiran scale
	desc, err := watermarkDesc("center", "50", true)
	if err != nil || !strings.HasPrefix(desc, "position:c") {
		t.Fatalf("desc tidak sesuai: %q err=%v", desc, err)
	}
}