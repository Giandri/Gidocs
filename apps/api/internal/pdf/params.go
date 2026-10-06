package pdf

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// ErrParam menandakan parameter user tidak lolos whitelist/validasi.
// Handler mengubahnya menjadi HTTP 400 dengan pesan yang aman ditampilkan.
var ErrParam = errors.New("parameter tidak valid")

var pageSelPattern = regexp.MustCompile(`^[0-9]{1,7}(-[0-9]{0,7})?(,[0-9]{1,7}(-[0-9]{0,7})?)*$`)

// pageSelection memvalidasi pilihan halaman seperti "1-3,5,8-" lalu memakai parser pdfcpu.
func pageSelection(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	if len(s) > 60 || !pageSelPattern.MatchString(s) {
		return nil, fmt.Errorf("%w: pakai angka seperti 1-3,5,8-", ErrParam)
	}
	return api.ParsePageSelection(s)
}

// checkPageBounds memastikan tiap token pilihan halaman tidak melewati jumlah halaman.
// Tanpa ini pdfcpu bisa menghasilkan PDF kosong untuk rentang di luar jangkauan.
func checkPageBounds(count int, sel []string) error {
	for _, token := range sel {
		from, thru, ok := splitRange(token)
		if !ok {
			return fmt.Errorf("%w: pakai angka seperti 1-3,5,8-", ErrParam)
		}
		if thru == 0 {
			thru = count // rentang terbuka "8-"
		}
		if from < 1 || from > count || thru < from || thru > count {
			return fmt.Errorf("%w: rentang melebihi jumlah halaman (%d)", ErrParam, count)
		}
	}
	return nil
}

// checkSelectionBounds mengambil jumlah halaman file lalu memvalidasi pilihan.
func checkSelectionBounds(ctx context.Context, in string, sel []string) error {
	count, err := api.PageCountFile(ctx, in)
	if err != nil {
		return errors.New("PDF tidak bisa dibaca")
	}
	return checkPageBounds(count, sel)
}

// splitRange memecah satu token "5", "3-7", atau "8-" (thru = 0 artinya terbuka).
func splitRange(token string) (from, thru int, ok bool) {
	sep := strings.IndexByte(token, '-')
	if sep < 0 {
		n, err := strconv.Atoi(token)
		return n, n, err == nil
	}
	from, err := strconv.Atoi(token[:sep])
	if err != nil {
		return 0, 0, false
	}
	rest := token[sep+1:]
	if rest == "" {
		return from, 0, true
	}
	thru, err = strconv.Atoi(rest)
	if err != nil {
		return 0, 0, false
	}
	return from, thru, true
}

func checkPassword(pw string) error {
	if pw == "" || len(pw) > 128 {
		return fmt.Errorf("%w: password wajib diisi (maks 128 karakter)", ErrParam)
	}
	return nil
}
