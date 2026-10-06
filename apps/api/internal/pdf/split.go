package pdf

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// maxSplitPages membatasi mode "every" agar satu PDF besar tidak membuat
// ribuan file sementara sekaligus.
const maxSplitPages = 500

// Split memecah PDF menjadi beberapa file di dalam outDir.
//
//	mode "every":  satu file per halaman (span 1).
//	mode "ranges": satu file per grup rentang pada parameter ranges ("1-3,5,8-").
func Split(ctx context.Context, in, outDir, mode, ranges string) error {
	switch mode {
	case "", "every":
		count, err := api.PageCountFile(ctx, in)
		if err != nil {
			return err
		}
		if count > maxSplitPages {
			return fmt.Errorf("%w: terlalu banyak halaman untuk dipecah per halaman (maks %d)", ErrParam, maxSplitPages)
		}
		return api.SplitFile(ctx, in, outDir, 1, nil)

	case "ranges":
		if strings.TrimSpace(ranges) == "" {
			return fmt.Errorf("%w: isi rentang halaman", ErrParam)
		}
		sel, err := pageSelection(ranges)
		if err != nil {
			return err
		}
		if err := checkSelectionBounds(ctx, in, sel); err != nil {
			return err
		}
		groups := strings.Split(ranges, ",")
		if len(groups) > 20 {
			return fmt.Errorf("%w: maksimal 20 rentang", ErrParam)
		}
		for i, group := range groups {
			group = strings.TrimSpace(group)
			out := filepath.Join(outDir, "part-"+strconv.Itoa(i+1)+".pdf")
			if err := api.TrimFile(ctx, in, out, []string{group}, nil); err != nil {
				return fmt.Errorf("gagal memotong rentang %q: %w", group, err)
			}
		}
		return nil

	default:
		return fmt.Errorf("%w: mode harus every atau ranges", ErrParam)
	}
}
