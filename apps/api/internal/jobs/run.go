package jobs

import (
	"context"
	"errors"
)

// Run menjalankan proses job sesuai tool pada payload.
// Mengembalikan error dengan pesan aman untuk ditampilkan ke user.
func Run(ctx context.Context, p Payload) error {
	switch p.Tool {
	case "docx-to-pdf", "excel-to-pdf", "pptx-to-pdf":
		return ConvertOffice(ctx, p)
	case "pdf-to-image":
		return ConvertToImages(ctx, p)
	case "ocr-pdf":
		return OCRPDF(ctx, p)
	default:
		return errors.New("tool tidak dikenal")
	}
}
