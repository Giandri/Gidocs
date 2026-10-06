"use client";

import { downloadBlob, type RunResult } from "../../lib/api";

interface ResultPanelProps {
  error: string | null;
  onReset: () => void;
  result: RunResult | null;
}

export function ResultPanel({ error, onReset, result }: ResultPanelProps) {
  if (error) {
    return (
      <p className="flex gap-2 text-[11px] leading-[1.8] text-error" role="alert">
        <span className="shrink-0 font-bold" aria-hidden="true">
          [!]
        </span>
        <span>{error}</span>
      </p>
    );
  }

  if (!result) return null;

  return (
    <div aria-live="polite">
      <p className="flex gap-2 text-[11px] leading-[1.8] text-ink">
        <span className="shrink-0 font-bold text-accent" aria-hidden="true">
          [ok]
        </span>
        <span>
          {result.filename} ({formatSize(result.blob.size)})
        </span>
      </p>
      <div className="mt-3 flex items-center gap-4">
        <button
          className="min-h-[36px] border border-accent px-4 text-[11px] text-accent transition-colors hover:bg-white hover:text-black focus-visible:bg-white focus-visible:text-black focus-visible:outline-none"
          onClick={() => downloadBlob(result)}
          type="button"
        >
          Download
        </button>
        <button
          className="min-h-[36px] px-1 text-[11px] text-nav transition-colors hover:text-white focus-visible:text-white focus-visible:outline-none"
          onClick={onReset}
          type="button"
        >
          Start over
        </button>
      </div>
    </div>
  );
}

function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
