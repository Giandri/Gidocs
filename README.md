# Gidocs

Open source PDF tools. No account, no database — files are uploaded, processed in a temporary folder, returned, then deleted.

## Run locally

```powershell
# requirements: bun, Go 1.24+, (optional) Ghostscript untuk Compress

bun install
turbo dev

# api saja
cd apps/api
go run ./cmd/api
```

## Tools

Merge, Split, Compress, Rotate, Watermark, Protect/Unlock, Remove/Reorder pages, Images to PDF — lihat daftar lengkap di halaman utama.

## License

**AGPL-3.0** (keputusan final). Wajib karena memakai Ghostscript (AGPL) untuk Compress dan rencana tool berat lainnya. Lihat `LICENSE`.
