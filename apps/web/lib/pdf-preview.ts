/**
 * Pratinjau PDF untuk Sign PDF: memuat dokumen dengan pdf.js dan merender
 * satu halaman ke canvas. Modul ini hanya dipakai halaman Sign PDF dan
 * diimpor secara dinamis supaya bundle tool lain tidak ikut memuat pdf.js.
 */

export interface PageGeometry {
  /** Jumlah halaman dokumen. */
  pageCount: number;
  /** Lebar halaman dalam poin PDF (1 pt = 1/72 inci). */
  widthPt: number;
  /** Tinggi halaman dalam poin PDF. */
  heightPt: number;
  /** Skala render piksel per poin. */
  scale: number;
  /** Lebar halaman hasil render dalam piksel. */
  widthPx: number;
  /** Tinggi halaman hasil render dalam piksel. */
  heightPx: number;
}

export interface PdfPreview {
  pageCount: number;
  /** Merender halaman 1-based ke canvas dan mengembalikan geometri hasil. */
  renderPage(pageNumber: number, canvas: HTMLCanvasElement): Promise<PageGeometry>;
  /** Merelease sumber daya dokumen. */
  destroy(): void;
}

type PdfjsLib = typeof import("pdfjs-dist");

let cached: Promise<PdfjsLib> | null = null;

// Skala render default: 1,4 px per pt ≈ 100 DPI, cukup tajam di layar 1x-2x.
const DEFAULT_SCALE = 1.4;

function loadPdfjs(): Promise<PdfjsLib> {
  cached ??= import("pdfjs-dist").then((lib) => {
    lib.GlobalWorkerOptions.workerSrc = new URL(
      "pdfjs-dist/build/pdf.worker.min.mjs",
      import.meta.url,
    ).toString();
    return lib;
  });
  return cached;
}

/** Memuat PDF dari ArrayBuffer (hasil File.arrayBuffer). */
export async function loadPdf(data: ArrayBuffer): Promise<PdfPreview> {
  const pdfjs = await loadPdfjs();
  const task = pdfjs.getDocument({ data: new Uint8Array(data) });
  const doc = await task.promise;

  return {
    pageCount: doc.numPages,
    async renderPage(pageNumber, canvas) {
      const page = await doc.getPage(pageNumber);
      const viewport = page.getViewport({ scale: DEFAULT_SCALE });
      const context = canvas.getContext("2d");
      if (!context) throw new Error("canvas tidak didukung browser ini");
      canvas.width = Math.round(viewport.width);
      canvas.height = Math.round(viewport.height);
      await page.render({ canvas, canvasContext: context, viewport }).promise;
      // Poin PDF ikut rotasi halaman yang sudah diterapkan viewport.
      return {
        pageCount: doc.numPages,
        widthPt: viewport.width / DEFAULT_SCALE,
        heightPt: viewport.height / DEFAULT_SCALE,
        scale: DEFAULT_SCALE,
        widthPx: canvas.width,
        heightPx: canvas.height,
      };
    },
    destroy() {
      void task.destroy();
    },
  };
}