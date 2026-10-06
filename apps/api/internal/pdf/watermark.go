package pdf

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// posWhitelist memetakan posisi UI ke parameter watermark pdfcpu.
var posWhitelist = map[string]string{
	"center":       "position:c",
	"top-left":     "position:tl",
	"bottom-right": "position:br",
	// diagonal: posisi tengah + diagonal menaik (2 = kiri atas ke kanan bawah)
	"diagonal": "position:c,diagonal:2",
}

// Watermark menambahkan teks watermark di semua halaman.
// opacityPersen: 10..100. position harus ada di posWhitelist.
func Watermark(ctx context.Context, in, out, text, position, opacityPersen string) error {
	text = strings.TrimSpace(text)
	if text == "" || len(text) > 100 {
		return fmt.Errorf("%w: teks watermark wajib diisi (maks 100 karakter)", ErrParam)
	}
	pos, ok := posWhitelist[position]
	if !ok {
		return fmt.Errorf("%w: posisi watermark tidak dikenal", ErrParam)
	}
	o, err := strconv.Atoi(opacityPersen)
	if err != nil || o < 10 || o > 100 {
		return fmt.Errorf("%w: opasitas harus 10-100", ErrParam)
	}
	desc := pos + ",opacity:" + strconv.FormatFloat(float64(o)/100, 'f', 2, 64)
	return api.AddTextWatermarksFile(ctx, in, out, nil, true, text, desc, nil)
}
