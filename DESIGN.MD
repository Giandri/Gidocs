# DESIGN: Gidocs

> Sistem desain Gidocs, diturunkan dari halaman beranda yang sudah ada (`apps/web`), plus spesifikasi halaman per tool yang belum dibuat.
> Aturan utama: **halaman baru harus terlihat seperti bagian dari beranda yang sama.** Jangan memperkenalkan gaya baru.

## 1. Karakter Visual

Bergaya **terminal / brutalist**: latar hitam, font monospace, garis tipis 1px, tanpa sudut membulat, tanpa bayangan, tanpa gradasi. Satu warna aksen kuning-oranye. Teks kecil dan padat. Interaksi ditandai dengan **inversi warna** (latar putih, teks hitam), bukan animasi.

Prinsip:
1. Garis, bukan kotak berwarna. Pemisah selalu `border` 1px.
2. Satu aksen (`#f3ab00`) untuk penanda, status, dan tautan.
3. Penanda teks seperti `[*]`, `[+]`, `[new]` menggantikan ikon.
4. Hover/focus = inversi (putih/hitam).
5. Tidak ada `rounded`, `shadow`, atau gradient.

## 2. Token

### Warna

| Token | Hex | Pemakaian |
|---|---|---|
| `bg` | `#000000` | Latar halaman |
| `line` | `#252525` | Semua border dan pemisah |
| `text` | `#d5d2d2` | Teks dasar |
| `text-strong` | `#f3f1f1` / `#e7e3e3` / `#e4e0e0` | Logo, judul, heading |
| `text-muted` | `#a29d9d` / `#969191` | Paragraf, deskripsi |
| `text-nav` | `#b8b4b4` | Tautan navigasi |
| `text-tool` | `#c7c2c2` | Label tombol tool |
| `text-faint` | `#797474` | Footer |
| `accent` | `#f3ab00` | Penanda, badge, tautan, logo "docs" |
| `accent-hover` | `#ffd36a` | Hover tautan aksen |
| `invert-bg` / `invert-fg` | `#ffffff` / `#000000` | State hover/focus |
| `error` (usulan) | `#ff6b6b` | Pesan error (belum ada di beranda) |

Saran implementasi: daftarkan sebagai variabel di Tailwind (`@theme` atau `tailwind.config`) agar tidak menyalin hex berulang.

### Tipografi

Font: `font-mono` untuk seluruh situs.

| Peran | Ukuran | Gaya |
|---|---|---|
| Logo | 20px | bold, `tracking-[-0.08em]`, "Gi" terang + "docs" aksen |
| H1 | 21px (24px di `sm`) | bold, `leading-[1.35]`, `tracking-[-0.045em]` |
| H2 (judul section) | 11px (12px di `sm`) | bold |
| Body / deskripsi | 10px (11px di `sm`) | `leading-[1.8]` sampai `1.9` |
| Label tool | 10px (11px di `sm`) | normal |
| Badge | 9px | warna aksen |
| Footer | 9px | `text-faint` |

### Layout dan Spasi

- Kontainer tengah: `max-w-[684px]`, `border-x` `line`, di dalam `main` hitam penuh layar.
- Konten teks di dalam section: `max-w-[584px]` (hero `max-w-[500px]`).
- Padding horizontal: `px-5` di mobile, `px-[30px]` (header) atau `px-[50px]` (section) di `sm`.
- Header: tinggi minimum `51px`, `border-b`.
- Antar section dipisah `border-b` atau `border-t` `line`.
- Breakpoint utama: `sm` (640px). Mobile memakai `MobileNav` dari `@repo/ui/mobile-nav`.

### Penanda Teks

| Penanda | Arti |
|---|---|
| `[*]` | Item langkah/daftar |
| `[+]` | Poin penting (privasi) |
| `[new]`, `[most]` | Badge tool, pojok kanan atas, aksen |
| `[!]` (usulan) | Error |
| `[ok]` (usulan) | Berhasil |

## 3. Komponen yang Sudah Ada (Beranda)

| Komponen | Spesifikasi |
|---|---|
| **Header/Nav** | Logo kiri (gambar `/gidocs.png` 16px + teks), tautan kanan (`GitHub`, `Docs`, `Stars`), disembunyikan di mobile diganti `MobileNav` |
| **Hero** | `min-h-[280px]`, rata kiri, H1 + satu paragraf |
| **How it works** | Daftar dengan penanda `[*]` aksen |
| **Grid tool** | `grid-cols-2`, `sm:grid-cols-5`, `border-l border-t` pada grid, `border-b border-r` pada tiap sel |
| **Sel tool** | `min-h-[89px]`, teks di tengah, hover/focus: `bg-white text-black border-white`, badge absolut `right-2 top-2` |
| **Privacy** | Heading + paragraf dengan penanda `[+]` dan tautan aksen bergaris bawah |
| **Footer** | `©2026 Giandri`, tengah, 9px |

## 4. Halaman Per Tool (Belum Dibuat)

### 4.1 Routing

Satu template, banyak tool: `apps/web/app/tools/[slug]/page.tsx`, dengan konfigurasi tool di satu file (`apps/web/lib/tools.ts`). Tombol di beranda diganti `Link` ke `/tools/<slug>`.

| Tool (label di beranda) | Slug | Input | Opsi | Endpoint | Status backend |
|---|---|---|---|---|---|
| Merge PDF | `merge-pdf` | banyak PDF | urutan file | `POST /api/merge` | Selesai |
| Split PDF | `split-pdf` | 1 PDF | per halaman / rentang | `POST /api/split` | Selesai |
| Compress PDF | `compress-pdf` | 1 PDF | level: light / medium / strong | `POST /api/compress` | Selesai (Ghostscript) |
| Watermark PDF | `watermark-pdf` | 1 PDF | teks, posisi, opasitas | `POST /api/watermark` | Selesai |
| Protect PDF | `protect-pdf` | 1 PDF | password | `POST /api/protect` | Selesai |
| JPEG to PDF | `jpeg-to-pdf` | banyak JPG/PNG | urutan, ukuran halaman | `POST /api/images-to-pdf` | Selesai |
| DOCX to PDF | `docx-to-pdf` | 1 DOCX | - | `POST /api/office-to-pdf` | Belum (LibreOffice, async) |
| Excel to PDF | `excel-to-pdf` | 1 XLSX | - | `POST /api/office-to-pdf` | Belum (LibreOffice, async) |
| PPTX to PDF | `pptx-to-pdf` | 1 PPTX | - | `POST /api/office-to-pdf` | Belum (LibreOffice, async) |
| Sign PDF | `sign-pdf` | 1 PDF | gambar/gambar tanda tangan, posisi | belum ditentukan | Belum, lihat catatan |

Catatan **Sign PDF**: kemungkinan paling aman dan privat dikerjakan **di browser** (pdf-lib + canvas tanda tangan). Putuskan di `PRD.md` sebelum dibuat.

Tool yang bergantung pada backend belum siap harus tampil sebagai **disabled** di beranda dengan badge `[soon]`, bukan mengarah ke halaman rusak.

### 4.2 Struktur Halaman

Memakai kerangka beranda yang sama (kontainer 684px, header, footer). Dari atas ke bawah:

```
[Header: sama dengan beranda]
-------------------------------------------------
< All tools                                 (tautan kecil, text-nav)

Merge PDF                                   (H1)
Combine multiple PDF files into one.        (deskripsi, text-muted)
-------------------------------------------------
[ Dropzone ]                                (border dashed, h-48)
-------------------------------------------------
Files (3)                                   (H2)
1. report.pdf          2.1 MB    [^] [v] [x]
2. appendix.pdf        0.4 MB    [^] [v] [x]
-------------------------------------------------
Options                                     (hanya jika tool punya opsi)
...
-------------------------------------------------
[ Run ]                                     (tombol utama)
-------------------------------------------------
Privacy note (1 baris, penanda [+])
-------------------------------------------------
[Footer]
```

### 4.3 Dropzone

- Kotak `border border-dashed border-[#252525]`, tinggi `h-48`, teks tengah 11px `text-muted`.
- Teks: `Drag PDF files here, or click to select`.
- State `dragging`: `border-[#f3ab00]`, teks aksen. Tanpa perubahan latar besar.
- Hover/focus: border putih (konsisten dengan sel tool).
- Seluruh area bisa diklik, dapat difokus keyboard (`role="button"`, `tabIndex=0`, Enter/Space membuka dialog file).
- Wajib `preventDefault()` pada `onDragOver` dan `onDrop`. Beri `pointer-events-none` pada elemen anak.

### 4.4 Daftar File

Baris dengan `border-b border-[#252525]`, teks 11px:

- Kiri: nomor + nama file (`truncate`), ukuran `text-faint`.
- Kanan: aksi teks `[^]` `[v]` `[x]`, hover inversi. Tombol `[^]` disabled pada baris pertama, `[v]` pada baris terakhir.
- Tombol aksi punya `aria-label` ("Move up", "Move down", "Remove").
- Untuk tool satu-file, tampilkan satu baris saja dengan aksi `[x]`.

### 4.5 Tombol Utama (Run)

- Lebar penuh, `min-h-[44px]`, 11px, `border border-[#252525]`.
- Aktif: `bg-white text-black` (inversi permanen sebagai CTA) atau `border-[#f3ab00] text-[#f3ab00]` dengan hover inversi. Pilih satu dan pakai konsisten di semua tool.
- Disabled: `opacity-40`, `cursor-not-allowed`, dengan teks alasan di bawahnya jika perlu (`Add at least 2 files`).
- Loading: teks `Processing...` dengan indikator teks sederhana (`[ ... ]` atau animasi titik). Tanpa spinner berwarna.

### 4.6 State Halaman

| State | Tampilan |
|---|---|
| **Idle** | Hanya dropzone, tombol Run disabled |
| **Files added** | Daftar file + opsi, Run aktif jika syarat terpenuhi |
| **Processing** | Dropzone dan daftar dikunci, Run berubah `Processing...` |
| **Done** | Blok hasil: `[ok] merged.pdf (3.4 MB)` + tombol `Download` + tautan `Start over` |
| **Error** | Baris `[!] <pesan dari backend>` warna `error`, file tetap ada agar user bisa mencoba lagi |

Pesan error ditampilkan apa adanya dari backend (sudah dirancang singkat dan aman). Tambahkan `role="alert"` agar terbaca pembaca layar.

### 4.7 Opsi Per Tool

| Tool | Kontrol |
|---|---|
| Split | Radio: `Every page` / `Page ranges` + input teks (contoh `1-3,5,8-`) |
| Compress | Radio tiga pilihan: `Light`, `Medium`, `Strong` (nilai whitelist yang sama dengan backend) |
| Watermark | Input teks, pilihan posisi, slider opasitas |
| Protect | Input password + konfirmasi password. Peringatan: `We cannot recover a lost password.` |
| JPEG to PDF | Pilihan ukuran halaman (`A4`, `Fit image`), urutan gambar |

Gaya kontrol: input `bg-black border border-[#252525]`, teks `#d5d2d2`, `focus` border putih. Label 10px `text-muted`. Tanpa sudut membulat.

### 4.8 Responsif

- Mobile: semua satu kolom, aksi baris file tetap terbaca (nama `truncate`, aksi tidak turun baris).
- Target sentuh tombol aksi minimal 32px tinggi walau teksnya kecil.
- Dropzone tetap bisa diklik di perangkat sentuh (tidak bergantung pada drag).

## 5. Aksesibilitas

- Kontras: `#d5d2d2` dan `#a29d9d` pada hitam sudah memadai. Hindari `#797474` untuk teks penting.
- Semua elemen interaktif punya `focus-visible` yang jelas (inversi putih).
- Dropzone dan tombol aksi dapat dioperasikan dengan keyboard.
- Error dan hasil diumumkan lewat `role="alert"` / `aria-live="polite"`.
- Jangan mengandalkan warna saja: error memakai `[!]`, sukses memakai `[ok]`.
- Gambar dekoratif memakai `alt=""` dan `aria-hidden`, seperti logo saat ini.

## 6. Struktur Kode yang Disarankan

```
apps/web/
├── app/
│   ├── page.tsx                  # beranda (sudah ada)
│   └── tools/[slug]/page.tsx     # template halaman tool
├── components/tools/
│   ├── tool-shell.tsx            # kerangka: header tool, judul, deskripsi, privacy note
│   ├── dropzone.tsx
│   ├── file-list.tsx
│   ├── run-button.tsx
│   └── result-panel.tsx
└── lib/
    ├── tools.ts                  # daftar tool: slug, judul, deskripsi, accept, opsi, endpoint, status
    └── api.ts                    # fungsi fetch ke /api/* + penanganan blob/unduhan
```

Komponen umum yang dipakai lintas halaman (`MobileNav`, header, footer) sebaiknya dipindahkan ke `packages/ui` agar beranda dan halaman tool memakai versi yang sama.

## 7. Catatan Perapian pada Beranda

Hal kecil yang ditemukan di `page.tsx`, sebaiknya dirapikan sekalian:

- Kelas `px-` pada `<main>` tidak valid (terpotong), hapus atau lengkapi.
- Variabel `i` pada `workflow.map` tidak dipakai.
- Nama produk tidak konsisten: "Gidocs" (logo) vs "Gidoc" (paragraf hero dan privacy). Pilih satu.
- Tautan `Stars` mengarah ke `#privacy`. Ganti ke halaman stargazers repo atau hapus.
- Tombol tool di beranda masih `type="button"` tanpa aksi. Ganti dengan `Link` ke `/tools/<slug>`.
- Kalimat privasi masih umum. Setelah backend berjalan, sebutkan fakta nyata: file diproses sementara dan dihapus setelah selesai.
- Padding berlebih di dalam section (`p-5`, `p-10` di dalam `px-5 py-5`) membuat ruang tidak seragam. Seragamkan.