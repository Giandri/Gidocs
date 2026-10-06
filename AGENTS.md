# AGENT: Panduan untuk AI Agent dan Kontributor

> Aturan kerja di repositori **Gidocs**. Baca `PRD.md` untuk tujuan produk dan `TASKS.md` untuk pekerjaan yang sedang berjalan. Jika ada konflik, `PRD.md` menang.

## 1. Gambaran Proyek

Gidocs adalah PDF tools open source tanpa akun dan tanpa database. File diunggah, diproses sementara, dikembalikan, lalu dihapus.

```
Gidocs/
├── apps/
│   ├── web/        # Next.js + TypeScript + Tailwind (frontend)
│   └── api/        # Go + chi + pdfcpu (backend)
│       ├── cmd/api/main.go
│       └── internal/
│           ├── http/     # handler, helper upload  (rencana: ganti nama ke handler)
│           ├── pdf/      # logika PDF (pdfcpu / CLI)
│           └── config/
├── packages/       # paket bersama (ui, config)
├── turbo.json
└── package.json
```

- Modul Go: `github.com/Giandri/Gidocs/apps/api`
- Package manager: **bun**
- Lingkungan utama pemilik proyek: **Windows / PowerShell**

## 2. Perintah Penting

```powershell
# dari root monorepo
turbo dev                      # jalankan web + api
turbo build
turbo test

# dari apps/api
go run ./cmd/api
go build ./...
go vet ./...
go test ./...
go mod tidy
govulncheck ./...              # cek kerentanan dependensi

# uji endpoint (PowerShell: pakai curl.exe)
curl.exe http://localhost:8080/api/health
curl.exe -X POST http://localhost:8080/api/merge -F "files=@a.pdf" -F "files=@b.pdf" -o hasil.pdf
```

Selalu jalankan `go build ./...` dan `go vet ./...` sebelum menyatakan perubahan Go selesai.

## 3. Aturan Backend (Go)

### Struktur dan gaya
- Satu tool = satu fungsi di `internal/pdf/` + satu handler + satu route di `main.go`.
- Fungsi di `internal/pdf/` **tidak** mengenal HTTP. Terima `context.Context` dan path file, kembalikan `error`.
- Handler tipis: validasi → simpan sementara → panggil fungsi → kirim hasil.
- Gunakan `log/slog`. Gunakan library standar dan chi. Jangan menambah framework baru tanpa alasan kuat.
- Import package `internal/http` memakai alias `apihttp` karena bentrok dengan `net/http`.
- Kode dan komentar boleh berbahasa Indonesia; pesan error ke user berbahasa Indonesia, singkat, dan aman ditampilkan.

### pdfcpu
- Versi terbaru mengharuskan `context.Context` sebagai argumen pertama pada banyak fungsi `api.*`.
- **Jangan menebak signature.** Periksa dengan:
  ```powershell
  go doc github.com/pdfcpu/pdfcpu/pkg/api NamaFungsi
  ```
- Kunci versi di `go.mod`.

### Pola handler wajib
1. `http.MaxBytesReader` untuk membatasi body.
2. Batasi jumlah file dan parameter.
3. Validasi magic bytes `%PDF-`, bukan ekstensi atau `Content-Type`.
4. `os.MkdirTemp` per request, `defer os.RemoveAll(dir)`.
5. Nama file sementara dibuat sendiri (`in-0.pdf`, `out.pdf`). **Nama dari user tidak pernah menjadi path.**
6. Teruskan `r.Context()` ke fungsi pemroses agar pembatalan dan timeout berlaku.
7. Kembalikan hasil dengan `Content-Type` dan `Content-Disposition` yang benar.

### CLI eksternal (Ghostscript, qpdf, Tesseract, LibreOffice)
- Selalu `exec.CommandContext` dengan timeout.
- Argumen **tidak boleh** disusun dari input user secara mentah. Gunakan whitelist untuk opsi (misalnya preset kompresi).
- Ghostscript wajib memakai `-dSAFER`.
- Jangan pernah memakai shell (`sh -c`, `cmd /c`) untuk merangkai perintah.
- Tangkap stderr untuk log internal, jangan tampilkan mentah ke user.

## 4. Aturan Frontend

- Next.js (App Router), TypeScript ketat, Tailwind.
- Komponen yang memakai state, event drag/drop, atau `fetch` ke browser harus berawalan `"use client"`.
- Semua panggilan backend lewat path relatif `/api/...` (dirute oleh rewrite Next.js).
- Kirim file dengan `FormData`. Nama field harus cocok dengan backend (`files`).
- Dropzone wajib memakai `preventDefault()` pada `onDragOver` dan `onDrop`.
- Tampilkan state loading, error dari backend, dan nonaktifkan tombol jika input belum memenuhi syarat.
- Validasi frontend hanya untuk kenyamanan. Validasi yang sah ada di Go.
- Teks UI harus melalui i18n (Indonesia dan Inggris), jangan hardcode jika sistem i18n sudah ada.

## 5. Keamanan (Tidak Bisa Ditawar)

Semua file dari user dianggap berbahaya.

- Jangan menyimpan file user di luar folder sementara per request.
- Jangan mencatat isi file atau nama file asli di log.
- Jangan menambah database, penyimpanan permanen, atau pelacak pihak ketiga tanpa revisi `PRD.md`.
- Jangan melonggarkan batas ukuran, jumlah file, atau timeout tanpa alasan tertulis.
- Jangan memakai library berlisensi komersial (misalnya unipdf).
- Jangan menambahkan secret ke repo. Gunakan `.env` (di-ignore) dan `.env.example`.
- Pesan error ke user tidak boleh membocorkan path server atau output mentah CLI.

## 6. Cara Menambah Tool Baru

1. Baca bagian tool terkait di `PRD.md` dan `TASKS.md`.
2. Buat fungsi di `apps/api/internal/pdf/<tool>.go` (terima `ctx`, path input, path output, parameter).
3. Buat handler di `internal/http/<tool>.go` mengikuti pola handler wajib (bagian 3).
4. Daftarkan route di `cmd/api/main.go`: `r.Post("/api/<tool>", apihttp.<Tool>Handler)`.
5. Buat unit test dengan PDF sampel kecil di `testdata/`.
6. Sambungkan UI di `apps/web` dan uji end-to-end dari browser.
7. Centang tugas di `TASKS.md`.

## 7. Pengujian

- Go: `go test ./...`. Uji minimal: kasus sukses, file bukan PDF, jumlah file salah, file korup, file terenkripsi.
- Frontend: uji manual drag/drop, klik pilih file, ubah urutan, hapus file, dan tampilan error.
- Setiap bug yang diperbaiki sebaiknya disertai test yang mencegah kemunculan ulang.

## 8. Masalah Umum

| Gejala | Penyebab | Solusi |
|---|---|---|
| `address already in use` di 8080 | Proses Go lama menggantung | `netstat -ano \| findstr :8080` lalu `taskkill /PID <PID> /F` |
| `not enough arguments in call to api.X` | Signature pdfcpu berubah (butuh `ctx`) | `go doc` untuk melihat signature |
| `apihttp undefined` | Import `internal/http` tanpa alias | Tambahkan alias `apihttp` |
| `go: cannot find main module` | Menjalankan `go` di luar `apps/api` | `cd apps/api` |
| `curl: (26)` | File lokal yang diunggah tidak ada | Pakai path file asli |
| Upload besar gagal (`413` / body exceeded) | Batas body proxy Next.js | Naikkan batas atau pakai CORS langsung ke Go |
| `air: command not found` | `GOPATH/bin` belum di PATH | Tambahkan ke PATH atau pakai `go run` |
| PowerShell: `mkdir -p` gagal | `-p` milik Linux | `mkdir a\b, c\d` |

## 9. Aturan Kerja untuk Agent

- Kerjakan **satu tugas dari `TASKS.md` per kali**, dan sebutkan tugas yang dikerjakan.
- Sebelum mengubah banyak file, jelaskan rencana singkat.
- Jangan menghapus atau menulis ulang file di luar lingkup tugas.
- Jangan mengarang API library. Jika ragu, periksa dokumentasi atau jalankan `go doc`.
- Setelah selesai: jalankan build dan test, lalu laporkan hasil dan langkah verifikasi yang bisa dicoba pemilik proyek.
- Jika instruksi bertentangan dengan `PRD.md` atau bagian Keamanan di atas, berhenti dan tanyakan.
- Berikan perintah yang kompatibel dengan **PowerShell** kecuali diminta lain.
- Perbarui `TASKS.md` dan dokumen ini jika arsitektur atau konvensi berubah.