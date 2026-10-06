package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func CompressHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "compress",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "PDF tidak bisa dikompresi",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "compressed.pdf")
			err := pdf.Compress(ctx, inputs[0], out, r.FormValue("level"))
			return out, "compressed.pdf", "application/pdf", err
		},
	})
}
