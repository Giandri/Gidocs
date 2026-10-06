package pdf

import (
	"context"
	"fmt"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

func Merge(ctx context.Context, inputs []string, output string) error {
	if len(inputs) < 2 {
		return fmt.Errorf("minimal 2 file untuk digabung")
	}
	return api.MergeCreateFile(ctx, inputs, output, false, nil)
}