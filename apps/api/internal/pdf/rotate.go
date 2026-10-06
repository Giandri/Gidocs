package pdf

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// Rotate memutar halaman PDF searah jarum jam. angle harus 90, 180, atau 270.
// pages kosong berarti semua halaman.
func Rotate(ctx context.Context, in, out, angle, pages string) error {
	a, err := strconv.Atoi(angle)
	if err != nil || (a != 90 && a != 180 && a != 270) {
		return fmt.Errorf("%w: sudut harus 90, 180, atau 270", ErrParam)
	}
	sel, err := pageSelection(pages)
	if err != nil {
		return err
	}
	if sel != nil {
		if err := checkSelectionBounds(ctx, in, sel); err != nil {
			return err
		}
	}
	return api.RotateFile(ctx, in, out, a, sel, nil)
}
