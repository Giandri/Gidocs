package pdf

import (
	"context"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Protect mengenkripsi PDF dengan password user (AES-256).
func Protect(ctx context.Context, in, out, password string) error {
	if err := checkPassword(password); err != nil {
		return err
	}
	conf := model.NewAESConfiguration(password, password, 256)
	return api.EncryptFile(ctx, in, out, conf)
}
