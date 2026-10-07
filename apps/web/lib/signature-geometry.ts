/**
 * Matematika koordinat untuk pratinjau Sign PDF.
 *
 * pdf.js merender halaman ke bitmap dengan skala tetap (px per pt), tetapi
 * canvas ditampilkan mengikuti lebar container (CSS `w-full`). Jadi ada dua
 * skala berbeda: skala render dan skala tampilan. Overlay tanda tangan harus
 * memakai skala tampilan, jika tidak posisinya meleset.
 */

/** Lebar halaman yang tampil di layar, dalam piksel CSS. */
export interface PageBox {
  width: number;
  height: number;
}

/** Skala tampilan (piksel CSS per titik PDF). Fallback 1 bila belum terukur. */
export function displayScale(cssWidth: number, widthPt: number): number {
  if (cssWidth <= 0 || widthPt <= 0) return 1;
  return cssWidth / widthPt;
}

/** Konversi titik PDF ke piksel CSS untuk penempatan overlay. */
export function ptToCss(pt: number, scale: number): number {
  return pt * scale;
}

/** Konversi piksel CSS (hasil drag) ke titik PDF bulat untuk dikirim ke backend. */
export function cssToPt(cssPx: number, scale: number): number {
  return Math.round(cssPx / scale);
}

/**
 * Menahan tanda tangan agar tidak keluar halaman. Margin dipakai untuk
 * menyisakan ruang kecil di tepi bawah.
 */
export function clampPosition(
  x: number,
  y: number,
  width: number,
  height: number,
  page: PageBox,
  margin = 24,
): { x: number; y: number } {
  return {
    x: Math.min(Math.max(0, x), Math.max(0, page.width - width)),
    y: Math.min(Math.max(0, y), Math.max(0, page.height - height - margin)),
  };
}
