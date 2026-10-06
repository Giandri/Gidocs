package pdf

import "github.com/pdfcpu/pdfcpu/pkg/api"

func init() {
	// Jangan tulis konfigurasi apa pun ke home user; semua kerja pakai folder sementara.
	api.DisableConfigDir()
}
