package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// RemovePages menghapus halaman terpilih dari PDF.
func RemovePages(ctx context.Context, in, out, pages string) error {
	sel, err := pageSelection(pages)
	if err != nil {
		return err
	}
	if len(sel) == 0 {
		return fmt.Errorf("%w: pilih halaman yang akan dihapus", ErrParam)
	}
	if err := checkSelectionBounds(ctx, in, sel); err != nil {
		return err
	}
	return api.RemovePagesFile(ctx, in, out, sel, nil)
}
