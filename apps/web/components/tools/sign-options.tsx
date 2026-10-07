"use client";

import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";

import { loadPdf, type PdfPreview } from "../../lib/pdf-preview";
import {
  clampPosition,
  cssToPt,
  displayScale,
  ptToCss,
  type PageBox,
} from "../../lib/signature-geometry";
import Checkbox from "../checkbox";
import WakeSlider from "../WakeSlider";

interface SignOptionsProps {
  /** PDF sumber untuk pratinjau. */
  source: File | null;
  fields: Record<string, string>;
  /** Perubahan posisi x, y (titik PDF) dan width. */
  onChange: (name: string, value: string) => void;
  /** Hasil tanda tangan PNG transparan. */
  onSignature: (file: File | null) => void;
}

const labelClass = "text-[10px] text-muted";
const optionClass = "flex cursor-pointer items-center gap-2 text-[11px] text-ink";
const SIGNATURE_FONT = '"Segoe Script", "Bradley Hand", "Brush Script MT", cursive';

// Batas lebar tanda tangan (pt); backend menerima 40-400.
const MIN_WIDTH_PT = 60;
const MAX_WIDTH_PT = 320;
const DEFAULT_WIDTH_PT = 150;
// Jarak dari tepi kiri dan tepi bawah untuk posisi tanda tangan awal.
const MARGIN_X_PT = 72;
const MARGIN_BOTTOM_PT = 144;

interface Geometry {
  widthPt: number;
  heightPt: number;
}

export function SignOptions({ source, fields, onChange, onSignature }: SignOptionsProps) {
  const [mode, setMode] = useState<"draw" | "type">("draw");
  const [name, setName] = useState("");
  const [signatureUrl, setSignatureUrl] = useState("");
  const [previewError, setPreviewError] = useState<string | null>(null);
  const [geometry, setGeometry] = useState<Geometry | null>(null);
  const [pageCount, setPageCount] = useState(0);
  const [resetToken, setResetToken] = useState(0);

  const pdfRef = useRef<PdfPreview | null>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const page = Number(fields.page || "1");

  // Muat PDF sumber setiap kali file berganti.
  useEffect(() => {
    if (!source) return;
    let cancelled = false;
    void (async () => {
      try {
        const pdf = await loadPdf(await source.arrayBuffer());
        if (cancelled) {
          pdf.destroy();
          return;
        }
        pdfRef.current?.destroy();
        pdfRef.current = pdf;
        setPageCount(pdf.pageCount);
        setPreviewError(null);
      } catch {
        if (!cancelled) setPreviewError("Pratinjau PDF gagal dimuat.");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [source]);

  // Render halaman aktif setiap kali halaman atau dokumen berubah.
  useEffect(() => {
    const pdf = pdfRef.current;
    const canvas = canvasRef.current;
    if (!pdf || !canvas) return;
    let cancelled = false;
    void pdf
      .renderPage(page, canvas)
      .then((result) => {
        if (cancelled) return;
        setGeometry({
          widthPt: result.widthPt,
          heightPt: result.heightPt,
        });
      })
      .catch(() => {
        if (!cancelled) setPreviewError("Halaman gagal dirender.");
      });
    return () => {
      cancelled = true;
    };
  }, [page, pageCount]);

  useEffect(() => () => pdfRef.current?.destroy(), []);

  // Lepaskan object URL lama agar tidak bocor.
  useEffect(() => {
    return () => {
      if (signatureUrl) URL.revokeObjectURL(signatureUrl);
    };
  }, [signatureUrl]);

  // Tanda tangan baru diletakkan di kanan bawah halaman.
  useEffect(() => {
    if (!signatureUrl || !geometry || fields.x) return;
    onChange("x", String(MARGIN_X_PT));
    onChange("y", String(Math.max(0, Math.round(geometry.heightPt - MARGIN_BOTTOM_PT))));
  }, [fields.x, geometry, onChange, signatureUrl]);

  const publish = useCallback(
    (file: File | null | Promise<File | null>) => {
      void Promise.resolve(file).then((result) => {
        if (signatureUrl) URL.revokeObjectURL(signatureUrl);
        setSignatureUrl(result ? URL.createObjectURL(result) : "");
        onSignature(result);
      });
    },
    [onSignature, signatureUrl],
  );

  return (
    <div className="space-y-3">
      <div className="flex flex-col gap-1">
        <span className={labelClass}>Signature</span>
        <Checkbox
          checked={mode === "draw"}
          className={optionClass}
          label="Draw"
          onCheckedChange={() => setMode("draw")}
        />
        <Checkbox
          checked={mode === "type"}
          className={optionClass}
          label="Type a name"
          onCheckedChange={() => setMode("type")}
        />
      </div>

      {mode === "draw" ? (
        <div className="relative">
          <SignaturePad
            resetToken={resetToken}
            onChange={publish}
          />
          <button
            className="absolute right-2 top-2 z-10 border border-line bg-black px-2 py-1 text-[11px] text-ink hover:border-white focus-visible:border-white focus-visible:outline-none"
            data-testid="signature-reset"
            onClick={() => {
              setResetToken((token) => token + 1);
              publish(null);
            }}
            type="button"
          >
            Reset
          </button>
        </div>
      ) : (
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Your name</span>
          <input
            className="border border-line bg-black px-2 py-1.5 text-[11px] text-ink focus:border-white focus:outline-none"
            onChange={(event) => {
              const value = event.target.value;
              setName(value);
              publish(value.trim() ? nameToPng(value.trim()) : null);
            }}
            placeholder="Full name"
            value={name}
          />
        </label>
      )}

      {source ? (
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className={labelClass}>Page</span>
            <span className="flex items-center gap-2 text-[11px] text-ink">
              <button
                className="border border-line px-2 py-0.5 disabled:opacity-40"
                disabled={page <= 1}
                onClick={() => onChange("page", String(page - 1))}
                type="button"
              >
                Prev
              </button>
              <span data-testid="sign-page">
                {page} / {Math.max(pageCount, 1)}
              </span>
              <button
                className="border border-line px-2 py-0.5 disabled:opacity-40"
                disabled={pageCount === 0 || page >= pageCount}
                onClick={() => onChange("page", String(page + 1))}
                type="button"
              >
                Next
              </button>
            </span>
          </div>
          <PdfPreviewArea
            canvasRef={canvasRef}
            error={previewError}
            fields={fields}
            geometry={geometry}
            onChange={onChange}
            signatureUrl={signatureUrl}
          />
        </div>
      ) : (
        <p className="text-[10px] text-muted">Add a PDF to see the preview.</p>
      )}

      <div className="space-y-0">
        <span className={labelClass}>Size ({fields.width || DEFAULT_WIDTH_PT} pt)</span>
        <WakeSlider
          ariaLabel="Signature size"
          defaultValue={DEFAULT_WIDTH_PT}
          min={MIN_WIDTH_PT}
          max={MAX_WIDTH_PT}
          step={5}
          bars={72}
          height={44}
          restHeight={10}
          gap={4}
          fillColor="#f3ab00"
          trackColor="#27272a"
          sensitivity={1}
          reach={6}
          skew={0.6}
          glide={0.3}
          smoothing={100}
          showValue
          onChange={(value) => onChange("width", String(value))}
          value={Number(fields.width || DEFAULT_WIDTH_PT)}
          disabled={false}
        />
      </div>
    </div>
  );
}

interface PadProps {
  resetToken: number;
  onChange: (file: File | null) => void;
}

/**
 * Kanvas gambar bebas. Coretan menumpuk: pengguna bisa menggambar berkali-kali
 * dan hanya perlu menekan Reset bila salah.
 */
function SignaturePad({ resetToken, onChange }: PadProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const drawingRef = useRef(false);

  // Siapkan ukuran kanvas + gaya goresan. MemAssign canvas.width menghapus
  // isi kanvas, jadi fungsi ini hanya mengubah ukuran bila benar-benar beda.
  const prepare = useCallback((clear: boolean) => {
    const canvas = canvasRef.current;
    if (!canvas) return null;
    const ctx = canvas.getContext("2d");
    if (!ctx) return null;
    const rect = canvas.getBoundingClientRect();
    const dpr = window.devicePixelRatio || 1;
    const width = Math.round(rect.width * dpr);
    const height = Math.round(rect.height * dpr);
    if (canvas.width !== width || canvas.height !== height) {
      canvas.width = width;
      canvas.height = height;
    }
    // setTransform dulu supaya skala tidak menumpuk saat dipanggil berulang.
    ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
    ctx.lineWidth = 2.5;
    ctx.lineCap = "round";
    ctx.lineJoin = "round";
    ctx.strokeStyle = "#000";
    if (clear) ctx.clearRect(0, 0, rect.width, rect.height);
    return { canvas, ctx, rect };
  }, []);

  // Kanvas dibersihkan hanya ketika tombol Reset ditekan.
  useEffect(() => {
    prepare(true);
  }, [prepare, resetToken]);

  // Ukuran kanvas mengikuti lebar container yang berubah.
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || typeof ResizeObserver === "undefined") return;
    const observer = new ResizeObserver(() => prepare(true));
    observer.observe(canvas);
    return () => observer.disconnect();
  }, [prepare]);

  const finish = async () => {
    if (!drawingRef.current) return;
    drawingRef.current = false;
    const canvas = canvasRef.current;
    if (!canvas) return;
    try {
      onChange(await trimmedPng(canvas));
    } catch {
      // Kanvas kosong: tidak ada tanda tangan.
    }
  };

  return (
    <canvas
      aria-label="Draw your signature"
      className="h-36 w-full touch-none cursor-crosshair border border-line bg-white"
      data-testid="signature-pad"
      onPointerDown={(event) => {
        const target = prepare(false);
        if (!target) return;
        target.canvas.setPointerCapture(event.pointerId);
        drawingRef.current = true;
        target.ctx.beginPath();
        target.ctx.moveTo(event.clientX - target.rect.left, event.clientY - target.rect.top);
        target.ctx.lineTo(event.clientX - target.rect.left + 0.1, event.clientY - target.rect.top);
        target.ctx.stroke();
      }}
      onPointerMove={(event) => {
        const canvas = canvasRef.current;
        if (!drawingRef.current || !canvas) return;
        const ctx = canvas.getContext("2d");
        if (!ctx) return;
        const rect = canvas.getBoundingClientRect();
        ctx.lineTo(event.clientX - rect.left, event.clientY - rect.top);
        ctx.stroke();
      }}
      onPointerUp={() => void finish()}
      ref={canvasRef}
    />
  );
}

interface PreviewAreaProps {
  canvasRef: React.RefObject<HTMLCanvasElement | null>;
  error: string | null;
  fields: Record<string, string>;
  geometry: Geometry | null;
  onChange: (name: string, value: string) => void;
  signatureUrl: string;
}

/** Pratinjau halaman PDF dengan tanda tangan yang bisa diseret. */
function PdfPreviewArea({ canvasRef, error, fields, geometry, onChange, signatureUrl }: PreviewAreaProps) {
  const [drag, setDrag] = useState<{ x: number; y: number } | null>(null);
  const [scale, setScale] = useState(1);
  const originRef = useRef<{ pointerX: number; pointerY: number; x: number; y: number } | null>(null);

  // Canvas dirender pada skala render pdf.js, tapi ditampilkan mengikuti lebar
  // container (CSS w-full). Overlay harus memakai skala tampilan yang diukur
  // dari lebar canvas sebenarnya, bukan skala render.
  useLayoutEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas || !geometry) return;
    const measure = () => {
      const width = canvas.getBoundingClientRect().width;
      if (width > 0) setScale(displayScale(width, geometry.widthPt));
    };
    measure();
    const observer = new ResizeObserver(measure);
    observer.observe(canvas);
    return () => observer.disconnect();
  }, [canvasRef, geometry]);

  const widthPt = Number(fields.width || DEFAULT_WIDTH_PT);
  const widthPx = geometry ? ptToCss(widthPt, scale) : 0;
  const page: PageBox = geometry
    ? { width: geometry.widthPt * scale, height: geometry.heightPt * scale }
    : { width: 0, height: 0 };
  const placed =
    drag ?? (geometry && fields.x ? { x: ptToCss(Number(fields.x), scale), y: ptToCss(Number(fields.y), scale) } : null);

  const clamp = (x: number, y: number, height: number) => clampPosition(x, y, widthPx, height, page);

  return (
    <div className="space-y-1">
      {error ? <p className="text-[10px] text-muted">{error}</p> : null}
      <div
        className="relative w-full touch-none select-none border border-line bg-black [-webkit-tap-highlight-color:transparent]"
        data-testid="sign-preview"
      >
        <canvas className="block h-auto w-full" ref={canvasRef} />
        {placed && geometry && signatureUrl ? (
          <img
            alt="Signature preview"
            className="absolute cursor-grab active:cursor-grabbing"
            data-testid="sign-overlay"
            draggable={false}
            onPointerDown={(event) => {
              event.currentTarget.setPointerCapture(event.pointerId);
              originRef.current = { pointerX: event.clientX, pointerY: event.clientY, x: placed.x, y: placed.y };
            }}
            onPointerMove={(event) => {
              const origin = originRef.current;
              if (!origin) return;
              setDrag(
                clamp(
                  origin.x + (event.clientX - origin.pointerX),
                  origin.y + (event.clientY - origin.pointerY),
                  event.currentTarget.offsetHeight,
                ),
              );
            }}
            onPointerUp={() => {
              const origin = originRef.current;
              if (origin && drag) {
                onChange("x", String(cssToPt(drag.x, scale)));
                onChange("y", String(cssToPt(drag.y, scale)));
              }
              originRef.current = null;
              setDrag(null);
            }}
            src={signatureUrl}
            style={{ left: placed.x, top: placed.y, width: widthPx }}
          />
        ) : null}
      </div>
      <p className="text-[10px] text-muted">Drag the signature to place it.</p>
    </div>
  );
}

// nameToPng merender nama menjadi PNG transparan bergaya tulisan tangan.
async function nameToPng(text: string): Promise<File> {
  const canvas = document.createElement("canvas");
  const ctx = canvas.getContext("2d");
  if (!ctx) throw new Error("canvas tidak didukung");
  const font = `72px ${SIGNATURE_FONT}`;
  ctx.font = font;
  canvas.width = Math.ceil(ctx.measureText(text).width) + 72;
  canvas.height = 144;
  ctx.font = font;
  ctx.fillStyle = "#000";
  ctx.textBaseline = "middle";
  ctx.fillText(text, 36, canvas.height / 2);
  return trimmedPng(canvas);
}

// trimmedPng memotong area kosong lalu mengembalikan PNG sebagai File.
async function trimmedPng(canvas: HTMLCanvasElement): Promise<File> {
  const ctx = canvas.getContext("2d");
  if (!ctx) throw new Error("canvas tidak didukung");
  const { data } = ctx.getImageData(0, 0, canvas.width, canvas.height);
  let minX = canvas.width;
  let minY = canvas.height;
  let maxX = -1;
  let maxY = -1;
  for (let y = 0; y < canvas.height; y++) {
    for (let x = 0; x < canvas.width; x++) {
      if ((data[(y * canvas.width + x) * 4 + 3] ?? 0) > 8) {
        minX = Math.min(minX, x);
        minY = Math.min(minY, y);
        maxX = Math.max(maxX, x);
        maxY = Math.max(maxY, y);
      }
    }
  }
  if (maxX < minX) throw new Error("tanda tangan kosong");

  const pad = 4;
  const sx = Math.max(0, minX - pad);
  const sy = Math.max(0, minY - pad);
  const cropped = document.createElement("canvas");
  cropped.width = Math.min(canvas.width, maxX + pad + 1) - sx;
  cropped.height = Math.min(canvas.height, maxY + pad + 1) - sy;
  const croppedCtx = cropped.getContext("2d");
  if (!croppedCtx) throw new Error("canvas tidak didukung");
  croppedCtx.drawImage(canvas, sx, sy, cropped.width, cropped.height, 0, 0, cropped.width, cropped.height);

  const blob = await new Promise<Blob>((resolve, reject) => {
    cropped.toBlob((result) => (result ? resolve(result) : reject(new Error("gagal membuat PNG"))), "image/png");
  });
  return new File([blob], "signature.png", { type: "image/png" });
}