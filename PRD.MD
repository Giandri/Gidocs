# PRD: Gidocs

> Dokumen kebutuhan produk. Sumber kebenaran untuk **apa** yang dibangun dan **kenapa**.
> Status: Draft v0.1

## 1. Ringkasan

**Gidocs** adalah website PDF tools open source yang bisa dipakai langsung tanpa akun. User membuka situs, memasukkan file, memproses, lalu mengunduh hasilnya. Tidak ada riwayat, tidak ada penyimpanan permanen.

## 2. Tujuan

1. Menyediakan tool PDF yang cepat, gratis, dan mudah dipakai (merge, split, rotate, compress, dan seterusnya).
2. Menjaga **privasi**: file hanya diproses sementara lalu dihapus.
3. Mudah di-**self-host** dengan satu perintah Docker.
4. Ramah kontributor open source: struktur jelas, pola penambahan tool seragam.
5. Mendukung **Bahasa Indonesia** dan Inggris.

## 3. Bukan Tujuan (Non-Goals)

- Akun user, login, riwayat file, atau penyimpanan cloud.
- Langganan berbayar dan billing.
- Editor PDF visual lengkap (edit teks bebas di dalam halaman).
- Kolaborasi real-time.

Jika salah satu dari ini dibutuhkan kelak, PRD ini harus direvisi lebih dulu.

## 4. Target Pengguna

| Pengguna                             | Kebutuhan                                         |
| ------------------------------------ | ------------------------------------------------- |
| Pelajar, mahasiswa, pekerja kantoran | Menggabung, memecah, memampatkan PDF dengan cepat |
| Orang yang peduli privasi            | Jaminan file tidak disimpan                       |
| Developer / komunitas                | Self-host dan kontribusi kode                     |

## 5. Prinsip Produk

1. **Stateless**: tanpa database. Status sementara (jika ada) hidup di memori atau Redis dengan TTL.
2. **Sekali pakai**: setiap request mandiri, file dihapus setelah selesai.
3. **Aman secara default**: semua file dianggap tidak tepercaya.
4. **Sederhana dulu**: proses sinkron untuk tool ringan, antrean hanya untuk tool berat.
5. **Library Go untuk operasi ringan, CLI eksternal untuk operasi berat.**

## 6. Arsitektur

```
[apps/web: Next.js]  --/api/*-->  [apps/api: Go + chi]  -->  [pdfcpu | Ghostscript | qpdf | Tesseract | LibreOffice]
                                        |
                                  Disk sementara (dihapus otomatis)
```

| Komponen                  | Teknologi                                          |
| ------------------------- | -------------------------------------------------- |
| Monorepo                  | Turborepo, package manager **bun**                 |
| Frontend                  | Next.js, TypeScript, Tailwind                      |
| Backend                   | Go 1.22+, chi, pdfcpu, `log/slog`                  |
| Tool berat (tahap lanjut) | Ghostscript, qpdf, Tesseract, LibreOffice headless |
| Antrean (tahap lanjut)    | Redis + asynq, hanya jika ada tool lama            |
| Deploy                    | Docker Compose + Caddy (reverse proxy, HTTPS)      |

Path modul Go: `github.com/Giandri/Gidocs/apps/api`.

## 7. Fitur

### 7.1 Prioritas

| Prioritas | Tool                          | Pendekatan                          | Sinkron/Async      |
| --------- | ----------------------------- | ----------------------------------- | ------------------ |
| P0        | Merge                         | pdfcpu                              | Sinkron            |
| P1        | Rotate                        | pdfcpu                              | Sinkron            |
| P1        | Split (hasil ZIP)             | pdfcpu                              | Sinkron            |
| P1        | Protect / Unlock              | pdfcpu                              | Sinkron            |
| P2        | Watermark                     | pdfcpu                              | Sinkron            |
| P2        | Hapus/ekstrak/urutkan halaman | pdfcpu                              | Sinkron            |
| P2        | Compress                      | Ghostscript                         | Sinkron atau async |
| P3        | PDF ke gambar                 | Poppler (`pdftoppm`)                | Async              |
| P3        | Gambar ke PDF                 | pdfcpu                              | Sinkron            |
| P4        | OCR                           | OCRmyPDF + Tesseract (`ind`, `eng`) | Async              |
| P4        | Word/Excel/PPT ke PDF         | LibreOffice headless                | Async              |

### 7.2 Kontrak API Umum

- Prefix: `/api/<tool>`, metode `POST`, `multipart/form-data`.
- Input file pada field `files` (atau `file` untuk satu file).
- Sukses: `200` dengan body biner hasil (`application/pdf` atau `application/zip`) dan `Content-Disposition: attachment`.
- Gagal: kode HTTP sesuai (`400`, `413`, `422`, `500`) dengan pesan teks singkat yang aman ditampilkan ke user.
- `GET /api/health` mengembalikan `{"status":"ok"}`.

## 8. Persyaratan Non-Fungsional

### Keamanan

- Validasi magic bytes `%PDF-`, jangan percaya ekstensi atau `Content-Type`.
- Nama file user **tidak pernah** dipakai sebagai path. Gunakan nama buatan sendiri.
- Batasi ukuran body (`http.MaxBytesReader`), jumlah file, dan timeout per request.
- Folder kerja sementara per request, selalu dihapus (`defer os.RemoveAll`).
- Argumen CLI eksternal tidak pernah disusun dari input mentah. Gunakan whitelist.
- Worker berat dijalankan non-root, filesystem read-only, tanpa akses jaringan keluar.
- Jangan mencatat isi file maupun nama file asli di log.
- Rate limiting di reverse proxy.

### Performa

- Tool ringan selesai di bawah 5 detik untuk file sampai 20 MB.
- Upload tidak dimuat penuh ke memori; gunakan streaming bila memungkinkan.

### Privasi

- Tidak ada analitik pihak ketiga yang melacak individu. Jika perlu, gunakan Plausible atau Umami.
- Tidak ada penyimpanan file setelah request selesai.

### Batas Awal

| Parameter               | Nilai   |
| ----------------------- | ------- |
| Ukuran total upload     | 50 MB   |
| Jumlah file per request | 20      |
| Timeout request         | 2 menit |

## 9. Lisensi

- Jika Ghostscript atau MuPDF dipakai, proyek dirilis dengan **AGPL-3.0**.
- Jika hanya pdfcpu (Apache 2.0) dan library permisif, lisensi bisa dipilih bebas, namun **AGPL-3.0** tetap disarankan agar konsisten ke depan.
- Hindari library dengan lisensi komersial (misalnya unipdf).

## 10. Metrik Keberhasilan

- Semua tool P0 sampai P2 berfungsi end-to-end dari UI.
- `docker compose up` menjalankan seluruh sistem tanpa konfigurasi tambahan.
- Tidak ada file user yang tersisa di disk setelah request selesai.
- Cakupan test untuk handler dan fungsi PDF inti.

## 11. Risiko

| Risiko                                      | Mitigasi                                                 |
| ------------------------------------------- | -------------------------------------------------------- |
| File berbahaya / PDF bomb                   | Batas ukuran, jumlah halaman, timeout, sandbox           |
| Kerentanan di Ghostscript/LibreOffice       | Jalankan di container terisolasi, perbarui image berkala |
| Batas body proxy Next.js untuk upload besar | Naikkan batasnya atau panggil Go langsung dengan CORS    |
| Biaya server membengkak                     | Rate limit, batas ukuran, auto-scaling bertahap          |
| Perubahan API pdfcpu antar versi            | Kunci versi di `go.mod`, cek signature dengan `go doc`   |

## 12. Pertanyaan Terbuka

- Apakah perlu halaman batas penggunaan dan kebijakan privasi sejak rilis pertama?
- Apakah tool berat (OCR, Office ke PDF) masuk rilis publik pertama atau menyusul?
- Apakah API akan dibuka untuk umum? Jika ya, itu mengubah keputusan soal database dan kuota.
