package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func RotateHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "rotate",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "PDF tidak bisa dirotasi",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "rotated.pdf")
			err := pdf.Rotate(ctx, inputs[0], out, r.FormValue("angle"), r.FormValue("pages"))
			return out, "rotated.pdf", "application/pdf", err
		},
	})
}
