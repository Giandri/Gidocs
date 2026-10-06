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

## Fase 1: Menyambungkan Merge ke Frontend (SELESAI)

- [x] Rewrite `/api/*` di `apps/web/next.config` ke `http://localhost:8080`
- [x] Verifikasi `turbo dev` dan `curl.exe http://localhost:3000/api/health`
- [x] Sambungkan dropzone + tombol Merge ke `POST /api/merge`
- [x] Urutan file bisa diubah (naik/turun) dan hapus per file
- [x] Tampilkan status loading dan pesan error dari backend
- [x] Uji manual: 2 PDF valid, 1 file saja, file palsu berekstensi `.pdf`, file terlalu besar
- [x] Putuskan solusi upload besar: CORS langsung ke Go (rewrite sudah bypass proxy Next.js)

## Fase 2: Merapikan Fondasi Backend (SELESAI)

- [x] Pindahkan `savePDF` ke `internal/handler/upload.go` (helper bersama)
- [x] Buat helper: `ParseUploads`, `NewWorkDir`, `WriteFileResponse`, `WriteError`
- [x] Ganti nama package `internal/http` menjadi `internal/handler` agar tidak bentrok dengan `net/http`
- [x] Konfigurasi lewat env: `PORT`, `MAX_UPLOAD_MB`, `MAX_FILES`, `REQUEST_TIMEOUT` (`internal/config`)
- [x] Tambah `apps/api/.env.example`
- [x] Middleware: request ID, logging (tanpa nama file), CORS
- [x] Middleware rate limit sederhana (atau serahkan ke Caddy) — diputuskan: delegasikan ke Caddy di Fase 7
- [x] Unit test untuk `pdf.Merge` dan validasi upload (`internal/pdf/merge_test.go`, `internal/handler/merge_test.go`)
- [x] Siapkan folder `testdata/` berisi PDF sampel kecil
- [x] Perbarui `.gitignore` (`apps/api/tmp/`, `apps/api/bin/`, `.env`)

## Fase 3: Tool Ringan (pdfcpu, sinkron) (SELESAI)

Setiap tool mengikuti pola: fungsi di `internal/pdf/` → handler → route → UI → test.
Catatan: Rotate, hapus halaman, reorder, dan unlock **API saja** (tidak ada di daftar beranda, lihat catatan di "Halaman Per Tool").

- [x] **Rotate**: `POST /api/rotate` (sudut 90/180/270, pilihan halaman)
- [x] **Split**: `POST /api/split` (per halaman atau rentang, hasil ZIP)
- [x] Hapus halaman: `POST /api/remove-pages`
- [x] Ekstrak/urutkan halaman: `POST /api/reorder`
- [x] **Protect**: `POST /api/protect` (password, AES-256)
- [x] **Unlock**: `POST /api/unlock` (password diberikan user)
- [x] Watermark teks: `POST /api/watermark`
- [x] Gambar ke PDF: `POST /api/images-to-pdf`
- [x] Validasi parameter dengan whitelist untuk semua opsi (angle, mode, rentang, posisi, opasitas, ukuran halaman, password; rentang juga dicek terhadap jumlah halaman karena pdfcpu bisa menulis PDF kosong)

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

---

## Tambahan: Halaman Per Tool (Frontend)

Rujukan desain: `DESIGN.md`. Daftar tool di beranda sedikit berbeda dari rencana awal (ada **Sign PDF** dan tiga tool **Office/JPEG ke PDF**, tidak ada **Rotate**). Selaraskan `PRD.md` dengan daftar ini.

- [x] Buat `apps/web/lib/tools.ts` (slug, judul, deskripsi, accept, opsi, endpoint, status)
- [x] Buat route `apps/web/app/tools/[slug]/page.tsx`
- [x] Komponen: `tool-shell`, `dropzone`, `file-list`, `run-button`, `result-panel`
- [x] `lib/api.ts`: fetch ke `/api/*`, penanganan blob dan unduhan
- [x] Ganti tombol tool di beranda menjadi `Link` ke `/tools/<slug>`
- [x] Tool yang backend-nya belum ada: tampil disabled dengan badge `[soon]`
- [x] Halaman **Merge PDF** tersambung ke `/api/merge` (tool pertama end-to-end)
- [x] Pindahkan header, footer, `MobileNav` ke `packages/ui`
- [x] Rapikan catatan di `DESIGN.md` bagian 7 (kelas `px-`, nama "Gidoc" vs "Gidocs", tautan Stars)
- [ ] Putuskan pendekatan **Sign PDF** (di browser atau backend) dan catat di `PRD.md`
- [ ] Tentukan urutan backend baru mengikuti beranda: Split, Compress, Watermark, Protect, JPEG to PDF, lalu Office to PDF

### Komponen bersama (ikuti `DESIGN.md` §4.3–§4.6)

- [x] `dropzone.tsx`: `border border-dashed`, `h-48`, teks tengah 11px, state `dragging` = border aksen, hover/focus = border putih, `preventDefault()` pada `onDragOver`/`onDrop`, anak `pointer-events-none`
- [x] Dropzone keyboard-operable: `role="button"`, `tabIndex=0`, Enter/Space buka dialog file, area penuh bisa diklik (termasuk di sentuh)
- [x] `file-list.tsx`: baris `border-b`, nomor + nama `truncate` + ukuran `text-faint`; aksi `[^] [v] [x]` hover inversi, `aria-label`, `[^]` disabled di baris pertama, `[v]` di baris terakhir
- [x] Tool satu-file hanya tampilkan satu baris dengan aksi `[x]`
- [x] `run-button.tsx`: lebar penuh, `min-h-[44px]`, 11px; pilih satu gaya aktif (border aksen + hover inversi) dan pakai konsisten di semua tool
- [x] Run disabled: `opacity-40`, `cursor-not-allowed` + teks alasan (mis. `Add at least 2 files`); loading: `Processing [ ... ]`, tanpa spinner berwarna
- [x] `result-panel.tsx`: blok sukses `[ok] nama (ukuran)` + `Download` + `Start over`; error `[!] <pesan backend>` warna error + `role="alert"`, file tetap tersedia untuk coba lagi
- [x] State halaman lengkap: idle (Run disabled), files added, processing (dropzone & daftar terkunci), done, error
- [x] `lib/api.ts`: fetch multipart ke `/api/*`, baca respons sebagai blob, unduh dengan nama dari `Content-Disposition`, pesan error dari body dipakai apa adanya

### Opsi per tool (`DESIGN.md` §4.7)

- [x] Split: radio `Every page` / `Page ranges` + input teks (`1-3,5,8-`)
- [x] Compress: radio `Light` / `Medium` / `Strong` (whitelist sama dengan backend)
- [x] Watermark: input teks, pilihan posisi, slider opasitas
- [x] Protect: password + konfirmasi password + peringatan `We cannot recover a lost password.`
- [x] JPEG to PDF: pilihan ukuran halaman (`A4`, `Fit image`)
- [x] Gaya kontrol seragam: input `bg-black border-line`, teks `#d5d2d2`, focus border putih, label 10px `text-muted`, tanpa rounded

### Aksesibilitas & responsif (`DESIGN.md` §4.8, §5)

- [x] `focus-visible` inversi putih di semua elemen interaktif
- [x] Hasil/error diumumkan `aria-live="polite"` / `role="alert"`; penanda `[!]`/`[ok]` sehingga tidak bergantung warna
- [x] Mobile satu kolom; aksi baris file tidak turun baris; target sentuh aksi minimal 32px
- [ ] Uji keyboard: tab ke dropzone, Enter/Space, urutkan/hapus file tanpa mouse

### Token & struktur (`DESIGN.md` §2, §6)

- [x] Daftarkan warna sebagai variabel Tailwind (`@theme`) agar hex tidak disalin berulang
- [x] Salin kerangka beranda (kontainer `max-w-[684px]`, header, footer) ke halaman tool agar terlihat satu desain
- [x] SEO dasar per halaman tool: `title` dan `description` dari `lib/tools.ts`

### Halaman per tool (end-to-end, urut)

- [x] Merge PDF → `POST /api/merge` (sudah ada backend)
- [x] Split PDF → `POST /api/split` (setelah backend Fase 3)
- [ ] Compress PDF → `POST /api/compress` (setelah Fase 4)
- [x] Watermark PDF → `POST /api/watermark`
- [x] Protect PDF → `POST /api/protect`
- [x] JPEG to PDF → `POST /api/images-to-pdf`
- [ ] DOCX/Excel/PPTX to PDF → `POST /api/office-to-pdf` (async, polling/SSE menyusul)
- [ ] Sign PDF → ikuti keputusan di `PRD.md`
- [x] Sebelum tiap backend siap: halaman tetap bisa dibuka, tombol Run disabled dengan alasan jelas (bukan halaman rusak)
