package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func ImagesToPDFHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   0,
		base:   "img",
		kind:   "gambar JPEG/PNG",
		magics: []uploadedKind{jpegKind, pngKind},
		fail:   "gambar tidak bisa diubah ke PDF",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "converted.pdf")
			err := pdf.ImagesToPDF(ctx, inputs, out, r.FormValue("size"))
			return out, "converted.pdf", "application/pdf", err
		},
	})
}
