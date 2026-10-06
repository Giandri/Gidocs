package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func ProtectHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "protect",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "PDF tidak bisa diproteksi",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "protected.pdf")
			err := pdf.Protect(ctx, inputs[0], out, r.FormValue("password"))
			return out, "protected.pdf", "application/pdf", err
		},
	})
}
