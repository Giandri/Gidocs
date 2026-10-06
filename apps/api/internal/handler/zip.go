package handler

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// zipDir mengarsipkan semua file .pdf di dalam dir ke zipPath.
func zipDir(dir, zipPath string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".pdf") {
			continue
		}
		src, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}
		w, err := zw.Create(entry.Name())
		if err != nil {
			src.Close()
			return err
		}
		if _, err := io.Copy(w, src); err != nil {
			src.Close()
			return err
		}
		src.Close()
	}
	return zw.Close()
}
