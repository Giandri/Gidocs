package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func UnlockHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "unlock",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "password salah atau file tidak terproteksi",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "unlocked.pdf")
			err := pdf.Unlock(ctx, inputs[0], out, r.FormValue("password"))
			return out, "unlocked.pdf", "application/pdf", err
		},
	})
}
