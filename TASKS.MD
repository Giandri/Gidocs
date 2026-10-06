# TASKS: Gidocs

> Daftar pekerjaan berurutan. Centang `[x]` saat selesai. Rujukan kebutuhan ada di `PRD.md`, aturan kerja di `AGENT.md`.

Legenda: `[x]` selesai, `[ ]` belum, `[~]` sedang dikerjakan.

---

## Fase 0: Fondasi (SELESAI)

- [x] Monorepo Turborepo dengan `apps/web` (Next.js)
- [x] Frontend UI selesai
- [x] Ganti isi `apps/api` dari Next.js ke Go
- [x] `go mod init github.com/Giandri/Gidocs/apps/api`
- [x] `cmd/api/main.go` dengan chi, graceful shutdown, `GET /api/health`
- [x] `package.json` pembungkus (`dev`, `build`, `test`, `lint`)
- [x] Endpoint `POST /api/merge` (validasi magic bytes, batas ukuran, folder sementara)

## Fase 1: Menyambungkan Merge ke Frontend

- [ ] Rewrite `/api/*` di `apps/web/next.config` ke `http://localhost:8080`
- [ ] Verifikasi `turbo dev` dan `curl.exe http://localhost:3000/api/health`
- [ ] Sambungkan dropzone + tombol Merge ke `POST /api/merge`
- [ ] Urutan file bisa diubah (naik/turun) dan hapus per file
- [ ] Tampilkan status loading dan pesan error dari backend
- [ ] Uji manual: 2 PDF valid, 1 file saja, file palsu berekstensi `.pdf`, file terlalu besar
- [ ] Putuskan solusi upload besar: naikkan batas body proxy atau CORS langsung ke Go

## Fase 2: Merapikan Fondasi Backend

- [ ] Pindahkan `savePDF` ke `internal/http/upload.go` (helper bersama)
- [ ] Buat helper: `parseUploads`, `newWorkDir`, `writeFileResponse`, `writeError`
- [ ] Ganti nama package `internal/http` menjadi `internal/handler` agar tidak bentrok dengan `net/http`
- [ ] Konfigurasi lewat env: `PORT`, `MAX_UPLOAD_MB`, `MAX_FILES`, `REQUEST_TIMEOUT`
- [ ] Tambah `apps/api/.env.example`
- [ ] Middleware: request ID, logging (tanpa nama file), CORS
- [ ] Middleware rate limit sederhana (atau serahkan ke Caddy)
- [ ] Unit test untuk `pdf.Merge` dan validasi upload
- [ ] Siapkan folder `testdata/` berisi PDF sampel kecil
- [ ] Perbarui `.gitignore` (`apps/api/tmp/`, `apps/api/bin/`, `.env`)

## Fase 3: Tool Ringan (pdfcpu, sinkron)

Setiap tool mengikuti pola: fungsi di `internal/pdf/` → handler → route → UI → test.

- [ ] **Rotate**: `POST /api/rotate` (sudut 90/180/270, pilihan halaman)
- [ ] **Split**: `POST /api/split` (per halaman atau rentang, hasil ZIP)
- [ ] **Hapus halaman**: `POST /api/remove-pages`
- [ ] **Ekstrak/urutkan halaman**: `POST /api/reorder`
- [ ] **Protect**: `POST /api/protect` (password)
- [ ] **Unlock**: `POST /api/unlock` (password diberikan user)
- [ ] **Watermark teks**: `POST /api/watermark`
- [ ] **Gambar ke PDF**: `POST /api/images-to-pdf`
- [ ] Validasi parameter dengan whitelist untuk semua opsi

## Fase 4: Compress (CLI eksternal pertama)

- [ ] Wrapper Ghostscript via `os/exec` dengan `-dSAFER`, timeout, dan `exec.CommandContext`
- [ ] Preset kompresi (whitelist): ringan, sedang, kuat
- [ ] `POST /api/compress` dan tampilkan ukuran sebelum/sesudah di UI
- [ ] Putuskan lisensi final (AGPL-3.0) dan catat di README

## Fase 5: Tool Berat (async)

- [ ] Pasang Redis + asynq, pisahkan `cmd/worker`
- [ ] Status job di Redis dengan TTL (tanpa database)
- [ ] Endpoint: `POST /api/jobs`, `GET /api/jobs/{id}`, `GET /api/jobs/{id}/download`
- [ ] PDF ke gambar (Poppler)
- [ ] OCR (OCRmyPDF + Tesseract bahasa `ind` dan `eng`)
- [ ] Office ke PDF (LibreOffice headless)
- [ ] Pembersih file kedaluwarsa terjadwal
- [ ] Polling/SSE di frontend

## Fase 6: Keamanan dan Hardening

- [ ] Batas jumlah halaman dan kompleksitas PDF (cegah PDF bomb)
- [ ] Worker non-root, filesystem read-only, tanpa jaringan keluar
- [ ] Batas CPU/RAM per container
- [ ] Audit log: pastikan tidak ada isi atau nama file asli
- [ ] Uji dengan PDF korup, terenkripsi, sangat besar, dan banyak halaman
- [ ] Tinjau dependensi (`govulncheck ./...`)

## Fase 7: Deployment

- [ ] `apps/api/Dockerfile` (multi-stage, distroless untuk API)
- [ ] Dockerfile worker (berisi Ghostscript, qpdf, Tesseract, LibreOffice)
- [ ] `docker-compose.yml`: web, api, (worker, redis), caddy
- [ ] Caddyfile: route `/api/*` ke Go, sisanya ke web, HTTPS otomatis
- [ ] GitHub Actions: lint, test Go dan web, build image
- [ ] Hosting: VPS minimal 2 vCPU / 4 GB RAM jika memakai tool berat

## Fase 8: Rilis Open Source

- [ ] `README.md` (cara jalan lokal, Docker, daftar tool)
- [ ] `LICENSE` (AGPL-3.0)
- [ ] `CONTRIBUTING.md` (cara menambah tool baru)
- [ ] `SECURITY.md` (cara melapor kerentanan)
- [ ] Halaman kebijakan privasi di frontend
- [ ] i18n Indonesia dan Inggris
- [ ] SEO dasar per halaman tool

---

## Catatan Teknis yang Perlu Diingat

- Lingkungan pengembangan: **Windows / PowerShell**. Gunakan `curl.exe`, bukan `curl`.
- Package manager: **bun**.
- Port: web `3000`, API `8080`. Hentikan server lama dengan `Ctrl+C` sebelum menjalankan `turbo dev`.
- Cek signature pdfcpu sesuai versi terpasang: `go doc github.com/pdfcpu/pdfcpu/pkg/api <NamaFungsi>`.
