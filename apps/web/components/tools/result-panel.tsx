"use client";

import { useState } from "react";

import { downloadBlob, type RunResult } from "../../lib/api";
import { canPreviewResult } from "../../lib/preview";
import { ResultPreview } from "./result-preview";

interface ResultPanelProps {
  error: string | null;
  /** Nama file unduhan hasil user (bawaan dari backend bila tidak diisi). */
  name?: string;
  onReset: () => void;
  result: RunResult | null;
}

export function ResultPanel({ error, name, onReset, result }: ResultPanelProps) {
  const [previewing, setPreviewing] = useState(false);

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
          {name || result.filename} ({formatSize(result.blob.size)})
        </span>
      </p>
      {result.beforeBytes !== undefined ? (
        <p className="mt-1 flex gap-2 text-[11px] leading-[1.8] text-nav">
          <span className="shrink-0 font-bold" aria-hidden="true">
            [*]
          </span>
          <span>
            {formatSize(result.beforeBytes)} -&gt; {formatSize(result.blob.size)} (
            {sizeDelta(result.beforeBytes, result.blob.size)})
          </span>
        </p>
      ) : null}
      <div className="mt-3 flex items-center gap-4">
        <button
          className="min-h-[36px] border border-accent px-4 text-[11px] text-accent transition-colors hover:bg-white hover:text-black focus-visible:bg-white focus-visible:text-black focus-visible:outline-none"
          onClick={() => downloadBlob(result, name)}
          type="button"
        >
          Download
        </button>
        {canPreviewResult(name || result.filename, result.blob.type) ? (
          <button
            aria-expanded={previewing}
            className="min-h-[36px] border border-line px-4 text-[11px] text-ink transition-colors hover:border-white focus-visible:border-white focus-visible:outline-none"
            data-testid="preview-toggle"
            onClick={() => setPreviewing((open) => !open)}
            type="button"
          >
            {previewing ? "Hide preview" : "Preview"}
          </button>
        ) : null}
        <button
          className="min-h-[36px] px-1 text-[11px] text-nav transition-colors hover:text-white focus-visible:text-white focus-visible:outline-none"
          onClick={onReset}
          type="button"
        >
          Start over
        </button>
      </div>
      {previewing ? <ResultPreview onClose={() => setPreviewing(false)} result={result} /> : null}
    </div>
  );
}

function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function sizeDelta(before: number, after: number): string {
  if (before <= 0) return "0%";
  const pct = Math.round((1 - after / before) * 100);
  return `${pct > 0 ? "-" : "+"}${Math.abs(pct)}%`;
}
