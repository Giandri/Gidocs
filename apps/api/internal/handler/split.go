package handler

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

func SplitHandler(w http.ResponseWriter, r *http.Request) {
	singleTool(w, r, toolSpec{
		want:   1,
		base:   "document",
		kind:   "PDF",
		magics: []uploadedKind{pdfKind},
		fail:   "PDF tidak bisa dibelah (cek rentang halaman)",
		fn: func(ctx context.Context, dir string, inputs []string) (string, string, string, error) {
			parts := filepath.Join(dir, "parts")
			if err := os.Mkdir(parts, 0o755); err != nil {
				return "", "", "", err
			}
			err := pdf.Split(ctx, inputs[0], parts, r.FormValue("mode"), r.FormValue("ranges"))
			if err != nil {
				return "", "", "", err
			}
			zipPath := filepath.Join(dir, "split.zip")
			if err := zipDir(parts, zipPath); err != nil {
				return "", "", "", err
			}
			return zipPath, "split.zip", "application/zip", nil
		},
	})
}
