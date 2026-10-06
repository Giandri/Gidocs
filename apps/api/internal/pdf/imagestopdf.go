package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

// ImagesToPDF mengubah gambar JPEG/PNG menjadi PDF.
//
//	size "a4": tiap gambar muat di halaman A4 (pusat, penuh).
//	size "fit": halaman dibuat seukuran gambar.
func ImagesToPDF(ctx context.Context, inputs []string, out, size string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("%w: minimal 1 gambar", ErrParam)
	}
	imp := pdfcpu.DefaultImportConfig()
	switch size {
	case "", "a4":
		imp.Pos = types.Center
		imp.Scale = 1
	case "fit":
		// default: halaman = ukuran gambar (Pos Full)
	default:
		return fmt.Errorf("%w: ukuran halaman harus a4 atau fit", ErrParam)
	}
	return api.ImportImagesFile(ctx, inputs, out, imp, nil)
}
