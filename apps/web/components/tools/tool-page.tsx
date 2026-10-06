"use client";

import { useState } from "react";

import { runTool, type RunResult } from "../../lib/api";
import { type OptionKind, type Tool } from "../../lib/tools";
import { Dropzone } from "./dropzone";
import { FileList } from "./file-list";
import { ResultPanel } from "./result-panel";
import { RunButton } from "./run-button";
import { ToolOptions } from "./tool-options";

const defaultFields: Record<string, Record<string, string>> = {
  none: {},
  split: { mode: "every", ranges: "" },
  compress: { level: "medium" },
  watermark: { text: "", position: "center", opacity: "50" },
  protect: { password: "", confirm: "" },
  "page-size": { size: "a4" },
};

export function ToolPage({ tool }: { tool: Tool }) {
  const [files, setFiles] = useState<File[]>([]);
  const [fields, setFields] = useState<Record<string, string>>(defaultFields[tool.options] ?? {});
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<RunResult | null>(null);

  const minFiles = tool.minFiles;
  const problem = validateFields(tool.options, fields);
  const ready = tool.ready && files.length >= minFiles && !problem && !busy;

  const reason = !tool.ready
    ? "This tool is not ready yet."
    : files.length < minFiles
      ? minFiles > 1
        ? `Add at least ${minFiles} files.`
        : "Add a file first."
      : problem;

  const clearOutput = () => {
    setError(null);
    setResult(null);
  };

  const addFiles = (incoming: File[]) => {
    const accepted = incoming.filter((file) => matchesAccept(file, tool.accept));
    if (accepted.length === 0) {
      setError("Unsupported file type.");
      return;
    }
    clearOutput();
    setFiles((previous) => {
      const next = tool.multiple ? [...previous, ...accepted] : accepted.slice(0, 1);
      return next.slice(0, tool.maxFiles);
    });
  };

  const moveFile = (index: number, delta: -1 | 1) => {
    clearOutput();
    setFiles((previous) => {
      const target = index + delta;
      if (target < 0 || target >= previous.length) return previous;
      const next = [...previous];
      const current = next[index];
      const neighbor = next[target];
      if (current === undefined || neighbor === undefined) return previous;
      next[index] = neighbor;
      next[target] = current;
      return next;
    });
  };

  const removeFile = (index: number) => {
    clearOutput();
    setFiles((previous) => previous.filter((_, i) => i !== index));
  };

  const changeField = (name: string, value: string) => {
    setFields((previous) => ({ ...previous, [name]: value }));
  };

  const run = async () => {
    setBusy(true);
    clearOutput();
    try {
      const output = await runTool(tool.endpoint, files, omitFields(fields, ["confirm"]));
      setResult(output);
    } catch (cause) {
      setError(cause instanceof Error ? cause.message : "Something went wrong.");
    } finally {
      setBusy(false);
    }
  };

  const reset = () => {
    clearOutput();
    setFiles([]);
    setFields(defaultFields[tool.options] ?? {});
  };

  return (
    <>
      <section className="border-b border-line px-5 py-6 sm:px-[50px]">
        <Dropzone
          accept={tool.accept}
          disabled={busy}
          hint={dropzoneHint(tool)}
          multiple={tool.multiple}
          onFiles={addFiles}
        />
      </section>

      {files.length > 0 ? (
        <section className="border-b border-line">
          <h2 className="border-b border-line px-5 py-2 text-[11px] font-bold text-[#e4e0e0] sm:px-[50px]">
            Files ({files.length})
          </h2>
          <FileList disabled={busy} files={files} onMove={moveFile} onRemove={removeFile} />
        </section>
      ) : null}

      {tool.options !== "none" ? (
        <section className="border-b border-line px-5 py-4 sm:px-[50px]">
          <h2 className="mb-3 text-[11px] font-bold text-[#e4e0e0]">Options</h2>
          <ToolOptions fields={fields} kind={tool.options} onChange={changeField} />
        </section>
      ) : null}

      <section className="border-b border-line px-5 py-4 sm:px-[50px]">
        <RunButton disabled={!ready} loading={busy} onClick={run} reason={reason} />
      </section>

      {error || result ? (
        <section className="border-b border-line px-5 py-4 sm:px-[50px]">
          <ResultPanel error={error} onReset={reset} result={result} />
        </section>
      ) : null}
    </>
  );
}

function validateFields(kind: OptionKind, fields: Record<string, string>): string | undefined {
  if (kind === "protect") {
    if (!fields.password) return "Enter a password.";
    if (fields.password !== fields.confirm) return "Passwords do not match.";
  }
  if (kind === "watermark" && !fields.text?.trim()) return "Enter watermark text.";
  if (kind === "split" && fields.mode === "ranges" && !fields.ranges?.trim()) return "Enter page ranges.";
  return undefined;
}

function matchesAccept(file: File, accept: string): boolean {
  const extensions = accept
    .split(",")
    .map((value) => value.trim().toLowerCase())
    .filter(Boolean);
  const name = file.name.toLowerCase();
  return extensions.some((extension) => name.endsWith(extension));
}

function dropzoneHint(tool: Tool): string {
  if (tool.accept.includes(".pdf")) return "Drag PDF files here, or click to select";
  if (tool.accept.includes(".jpg") || tool.accept.includes(".png")) {
    return "Drag image files here, or click to select";
  }
  return "Drag your file here, or click to select";
}

function omitFields(fields: Record<string, string>, names: string[]): Record<string, string> {
  return Object.fromEntries(Object.entries(fields).filter(([name]) => !names.includes(name)));
}
