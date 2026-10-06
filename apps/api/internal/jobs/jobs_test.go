package jobs

import (
	"archive/zip"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidID(t *testing.T) {
	valid := NewID()
	if len(valid) != 32 || !ValidID(valid) {
		t.Fatalf("NewID tidak valid: %q", valid)
	}
	for _, bad := range []string{
		"", "abc", "../../../etc/passwd", "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"0123456789abcdef0123456789abcde", // 31 char
	} {
		if ValidID(bad) {
			t.Errorf("ValidID(%q) = true, ingin false", bad)
		}
	}
}

func TestOfficeExt(t *testing.T) {
	cases := map[string]string{
		"docx-to-pdf":  ".docx",
		"excel-to-pdf": ".xlsx",
		"pptx-to-pdf":  ".pptx",
	}
	for slug, want := range cases {
		got, ok := OfficeExt(slug)
		if !ok || got != want {
			t.Errorf("OfficeExt(%q) = %q,%v; ingin %q,true", slug, got, ok, want)
		}
	}
	for _, slug := range []string{"", "merge-pdf", "../etc", "office-to-pdf"} {
		if _, ok := OfficeExt(slug); ok {
			t.Errorf("OfficeExt(%q) lolos whitelist", slug)
		}
	}
}

func TestSweep(t *testing.T) {
	root := t.TempDir()
	old := filepath.Join(root, NewID())
	fresh := filepath.Join(root, NewID())
	if err := os.MkdirAll(old, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(fresh, 0o700); err != nil {
		t.Fatal(err)
	}
	bukanJob := filepath.Join(root, "bukan-job.txt")
	if err := os.WriteFile(bukanJob, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	lama := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(old, lama, lama); err != nil {
		t.Fatal(err)
	}

	removed, err := Sweep(root, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Fatalf("removed = %d, ingin 1", removed)
	}
	if _, err := os.Stat(old); !os.IsNotExist(err) {
		t.Error("folder lama seharusnya terhapus")
	}
	if _, err := os.Stat(fresh); err != nil {
		t.Error("folder baru seharusnya dipertahankan")
	}
	if _, err := os.Stat(bukanJob); err != nil {
		t.Error("file non-job seharusnya dipertahankan")
	}
}

// writeMinimalDOCX membuat DOCX minimal (zip) untuk uji konversi.
func writeMinimalDOCX(t *testing.T, path string) {
	t.Helper()
	files := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		"word/document.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
<w:body><w:p><w:r><w:t>Uji konversi.</w:t></w:r></w:p></w:body></w:document>`,
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestConvertOffice(t *testing.T) {
	if _, err := FindSOFFICE(); err != nil {
		t.Skip("LibreOffice tidak tersedia, tes dilewati")
	}
	id := NewID()
	dir := JobDir(id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	input := filepath.Join(dir, "input.docx")
	writeMinimalDOCX(t, input)

	payload := Payload{JobID: id, Tool: "docx-to-pdf", Input: input, OutDir: dir}
	if err := ConvertOffice(context.Background(), payload); err != nil {
		t.Fatalf("ConvertOffice: %v", err)
	}
	data, err := os.ReadFile(ResultPath(id))
	if err != nil {
		t.Fatalf("hasil tidak ada: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("hasil bukan PDF: %.8q", data[:min(8, len(data))])
	}
}
