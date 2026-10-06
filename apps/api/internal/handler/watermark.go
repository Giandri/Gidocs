package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "watermark",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "watermark tidak bisa ditambahkan",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "watermarked.pdf")
			err := pdf.Watermark(ctx, inputs[0], out,
				r.FormValue("text"), r.FormValue("position"), r.FormValue("opacity"))
			return out, "watermarked.pdf", "application/pdf", err
		},
	})
}
