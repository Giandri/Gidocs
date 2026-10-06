package handler

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func RemovePagesHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "remove",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "halaman tidak bisa dihapus (cek pilihan halaman)",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			out := filepath.Join(dir, "removed.pdf")
			err := pdf.RemovePages(ctx, inputs[0], out, r.FormValue("pages"))
			return out, "removed.pdf", "application/pdf", err
		},
	})
}
