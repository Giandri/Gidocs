package pdf

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

// posWhitelist memetakan posisi UI ke parameter watermark pdfcpu.
var posWhitelist = map[string]string{
	"center":       "position:c",
	"top-left":     "position:tl",
	"bottom-right": "position:br",
	// diagonal: posisi tengah + diagonal menaik (2 = kiri atas ke kanan bawah)
	"diagonal": "position:c,diagonal:2",
}

// watermarkDesc membangun deskripsi posisi+opasitas untuk pdfcpu.
// diagonal hanya valid untuk watermark teks (parse lain untuk gambar).
func watermarkDesc(position, opacityPersen string, allowDiagonal bool) (string, error) {
	pos, ok := posWhitelist[position]
	if !ok {
		return "", fmt.Errorf("%w: posisi watermark tidak dikenal", ErrParam)
	}
	if !allowDiagonal && position == "diagonal" {
		return "", fmt.Errorf("%w: posisi diagonal hanya untuk watermark teks", ErrParam)
	}
	o, err := strconv.Atoi(opacityPersen)
	if err != nil || o < 10 || o > 100 {
		return "", fmt.Errorf("%w: opasitas harus 10-100", ErrParam)
	}
	return pos + ",opacity:" + strconv.FormatFloat(float64(o)/100, 'f', 2, 64), nil
}

// sizeFactor mengubah persen ukuran (10-100) menjadi faktor skala 0.10-1.00.
// Kosong dianggap 50 (nilai default).
func sizeFactor(sizePersen string) (float64, error) {
	if strings.TrimSpace(sizePersen) == "" {
		sizePersen = "50"
	}
	s, err := strconv.Atoi(sizePersen)
	if err != nil || s < 10 || s > 100 {
		return 0, fmt.Errorf("%w: ukuran harus 10-100", ErrParam)
	}
	return float64(s) / 100, nil
}

func Watermark(ctx context.Context, in, out, text, position, opacityPersen, sizePersen string) error {
	text = strings.TrimSpace(text)
	if text == "" || len(text) > 100 {
		return fmt.Errorf("%w: teks watermark wajib diisi (maks 100 karakter)", ErrParam)
	}
	desc, err := watermarkDesc(position, opacityPersen, true)
	if err != nil {
		return err
	}
	f, err := sizeFactor(sizePersen)
	if err != nil {
		return err
	}
	// Ukuran teks = persen ukuran dikali 120 -> 50% = 60pt (default pdfcpu).
	desc += ",points:" + strconv.Itoa(int(math.Round(f*120)))
	return api.AddTextWatermarksFile(ctx, in, out, nil, true, text, desc, nil)
}

// WatermarkImage menempatkan gambar (PNG/JPEG) sebagai watermark di semua
// halaman. Ukuran relatif lebar halaman sesuai param size (10-100%, default 50).
func WatermarkImage(ctx context.Context, in, out, imgPath, position, opacityPersen, sizePersen string) error {
	desc, err := watermarkDesc(position, opacityPersen, false)
	if err != nil {
		return err
	}
	f, err := sizeFactor(sizePersen)
	if err != nil {
		return err
	}
	scale := strconv.FormatFloat(f, 'f', 2, 64)
	return api.AddImageWatermarksFile(ctx, in, out, nil, true, imgPath, desc+",scale:"+scale, nil)
}
