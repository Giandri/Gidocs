package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func ReorderHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "reorder",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "halaman tidak bisa diurutkan (cek urutan halaman)",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "reordered.pdf")
			err := pdf.Reorder(ctx, inputs[0], out, r.FormValue("order"))
			return out, "reordered.pdf", "application/pdf", err
		},
	})
}
