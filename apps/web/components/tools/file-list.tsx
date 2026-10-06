"use client";

import { type ReactNode } from "react";

interface FileListProps {
  files: File[];
  disabled?: boolean;
  onMove: (index: number, delta: -1 | 1) => void;
  onRemove: (index: number) => void;
}

function RowAction({
  label,
  disabled,
  onClick,
  children,
}: {
  label: string;
  disabled: boolean;
  onClick: () => void;
  children: ReactNode;
}) {
  return (
    <button
      aria-label={label}
      className="flex min-h-[32px] min-w-[32px] items-center justify-center px-1 text-[11px] text-nav transition-colors hover:bg-white hover:text-black focus-visible:bg-white focus-visible:text-black focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-transparent disabled:hover:text-nav"
      disabled={disabled}
      onClick={onClick}
      type="button"
    >
      {children}
    </button>
  );
}

export function FileList({ files, disabled = false, onMove, onRemove }: FileListProps) {
  return (
    <ul>
      {files.map((file, index) => (
        <li
          className="flex items-center gap-3 border-b border-line px-5 py-1 text-[11px] sm:px-[50px]"
          key={`${file.name}-${index}`}
        >
          <span className="min-w-0 flex-1 truncate text-ink">
            {index + 1}. {file.name}
          </span>
          <span className="shrink-0 text-[9px] text-faint">{formatSize(file.size)}</span>
          <span className="flex shrink-0 gap-1">
            <RowAction disabled={disabled || index === 0} label="Move up" onClick={() => onMove(index, -1)}>
              [^]
            </RowAction>
            <RowAction
              disabled={disabled || index === files.length - 1}
              label="Move down"
              onClick={() => onMove(index, 1)}
            >
              [v]
            </RowAction>
            <RowAction disabled={disabled} label="Remove" onClick={() => onRemove(index)}>
              [x]
            </RowAction>
          </span>
        </li>
      ))}
    </ul>
  );
}

function formatSize(bytes: number): string {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}
