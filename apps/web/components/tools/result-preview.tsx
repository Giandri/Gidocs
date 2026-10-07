"use client";

import { useEffect, useRef, useState } from "react";

import { type RunResult } from "../../lib/api";
import { loadPdf, type PdfPreview } from "../../lib/pdf-preview";

interface ResultPreviewProps {
  result: RunResult;
  onClose: () => void;
}

/** Pratinjau hasil PDF di dalam halaman supaya user tak perlu mengunduh dulu. */
export function ResultPreview({ result, onClose }: ResultPreviewProps) {
  const [pdf, setPdf] = useState<PdfPreview | null>(null);
  const [page, setPage] = useState(1);
  const [pageCount, setPageCount] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // Muat PDF hasil saat pratinjau dibuka. Blob hanya dibaca di browser.
  useEffect(() => {
    let cancelled = false;
    let loaded: PdfPreview | null = null;
    void (async () => {
      try {
        const buffer = await result.blob.arrayBuffer();
        const preview = await loadPdf(buffer);
        loaded = preview;
        if (cancelled) {
          preview.destroy();
          return;
        }
        setPdf((previous) => {
          previous?.destroy();
          return preview;
        });
        setPageCount(preview.pageCount);
        setPage(1);
        setError(null);
      } catch {
        if (!cancelled) setError("Pratinjau gagal dimuat.");
      }
    })();
    return () => {
      cancelled = true;
      loaded?.destroy();
    };
  }, [result]);

  useEffect(() => () => pdf?.destroy(), [pdf]);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!pdf || !canvas) return;
    let cancelled = false;
    void pdf.renderPage(page, canvas).catch(() => {
      if (!cancelled) setError("Halaman gagal dirender.");
    });
    return () => {
      cancelled = true;
    };
  }, [page, pdf]);

  return (
    <div className="mt-3 space-y-2 border border-line p-3" data-testid="result-preview">
      <div className="flex items-center justify-between">
        <span className="text-[10px] text-muted">Preview</span>
        <span className="flex items-center gap-2 text-[11px] text-ink">
          <button
            className="border border-line px-2 py-0.5 disabled:opacity-40"
            disabled={page <= 1}
            onClick={() => setPage((current) => Math.max(1, current - 1))}
            type="button"
          >
            Prev
          </button>
          <span data-testid="preview-page">
            {page} / {Math.max(pageCount, 1)}
          </span>
          <button
            className="border border-line px-2 py-0.5 disabled:opacity-40"
            disabled={pageCount === 0 || page >= pageCount}
            onClick={() => setPage((current) => Math.min(pageCount, current + 1))}
            type="button"
          >
            Next
          </button>
          <button className="px-1 text-nav hover:text-white" onClick={onClose} type="button">
            Close
          </button>
        </span>
      </div>
      {error ? <p className="text-[10px] text-muted">{error}</p> : null}
      <canvas className="block h-auto w-full border border-line" ref={canvasRef} />
    </div>
  );
}
