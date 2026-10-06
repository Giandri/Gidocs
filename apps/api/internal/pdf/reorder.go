package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Reorder mengurutkan/mengekstrak halaman sesuai pilihan user, misal "3,1,2".
func Reorder(ctx context.Context, in, out, order string) error {
	sel, err := pageSelection(order)
	if err != nil {
		return err
	}
	if len(sel) == 0 {
		return fmt.Errorf("%w: isi urutan halaman", ErrParam)
	}
	if err := checkSelectionBounds(ctx, in, sel); err != nil {
		return err
	}
	return api.CollectFile(ctx, in, out, sel, nil)
}
