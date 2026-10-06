package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Giandri/Gidocs/apps/api/internal/pdf"
)

// toolSpec menjelaskan handler berformat satu upload: file yang diterima dan cara memprosesnya.
type toolSpec struct {
	want   int      // jumlah file wajib; 0 = satu file atau lebih
	base   string   // nama dasar file sementara (bukan nama dari user)
	kind   string   // label jenis file untuk pesan validasi
	magics []uploadedKind
	fail   string // pesan gagal proses, aman ditampilkan ke user
	fn     func(ctx context.Context, dir string, inputs []string) (out, filename, contentType string, err error)
}

// singleTool adalah kerangka handler wajib: validasi upload -> simpan sementara -> proses -> kirim hasil.
func singleTool(w http.ResponseWriter, r *http.Request, spec toolSpec) {
	cfg := LoadUploadConfig()
	files, err := ParseUploads(r, "files", cfg)
	if err != nil {
		switch err {
		case ErrUploadTooLarge:
			WriteError(w, "upload terlalu besar atau tidak valid", http.StatusRequestEntityTooLarge)
		case ErrNoFiles:
			WriteError(w, "kirim file pada field 'files'", http.StatusBadRequest)
		case ErrTooManyFiles:
			WriteError(w, "terlalu banyak file (maks "+strconv.Itoa(cfg.MaxFiles)+")", http.StatusBadRequest)
		default:
			WriteError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if spec.want > 0 && len(files) != spec.want {
		WriteError(w, "kirim tepat "+strconv.Itoa(spec.want)+" file", http.StatusBadRequest)
		return
	}

	dir, err := NewWorkDir(spec.base + "-*")
	if err != nil {
		WriteError(w, "gagal menyiapkan folder kerja", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(dir)

	inputs := make([]string, 0, len(files))
	for i, fh := range files {
		index := i
		if spec.want == 1 {
			index = -1
		}
		path, err := saveUpload(fh, dir, spec.base, index, spec.magics, spec.kind)
		if err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		inputs = append(inputs, path)
	}

	out, name, contentType, err := spec.fn(r.Context(), dir, inputs)
	if err != nil {
		if errors.Is(err, pdf.ErrParam) {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Log untuk debugging internal; pesan ke user tetap generik.
		slog.Error("pemrosesan tool gagal", "tool", spec.base, "err", err)
		WriteError(w, spec.fail, http.StatusUnprocessableEntity)
		return
	}
	if err := WriteFileResponse(w, out, name, contentType); err != nil {
		WriteError(w, "gagal mengirim hasil", http.StatusInternalServerError)
	}
}
