package pdf

import (
	"context"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Unlock membuka enkripsi PDF memakai password yang diberikan user.
func Unlock(ctx context.Context, in, out, password string) error {
	if err := checkPassword(password); err != nil {
		return err
	}
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password
	return api.DecryptFile(ctx, in, out, conf)
}
