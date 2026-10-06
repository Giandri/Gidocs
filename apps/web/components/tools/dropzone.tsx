"use client";

import { type ChangeEvent, type DragEvent, type KeyboardEvent, useRef, useState } from "react";

interface DropzoneProps {
  accept: string;
  hint: string;
  multiple: boolean;
  disabled?: boolean;
  onFiles: (files: File[]) => void;
}

export function Dropzone({ accept, hint, multiple, disabled = false, onFiles }: DropzoneProps) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [dragging, setDragging] = useState(false);

  const pick = (event: ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) onFiles(Array.from(event.target.files));
    event.target.value = "";
  };

  const open = () => {
    if (!disabled) inputRef.current?.click();
  };

  const onDragOver = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    if (!disabled) setDragging(true);
  };

  const onDrop = (event: DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    setDragging(false);
    if (disabled) return;
    onFiles(Array.from(event.dataTransfer.files));
  };

  const onKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    if (event.key === "Enter" || event.key === " ") {
      event.preventDefault();
      open();
    }
  };

  const state = dragging
    ? "border-accent text-accent"
    : disabled
      ? "border-line text-muted opacity-40"
      : "border-line text-muted hover:border-white focus-visible:border-white";

  return (
    <div
      aria-disabled={disabled}
      className={`flex h-48 w-full cursor-pointer flex-col items-center justify-center border border-dashed px-4 text-center focus-visible:outline-none ${state}`}
      onClick={open}
      onKeyDown={onKeyDown}
      onDragLeave={() => setDragging(false)}
      onDragOver={onDragOver}
      onDrop={onDrop}
      role="button"
      tabIndex={disabled ? -1 : 0}
    >
      <input accept={accept} className="hidden" multiple={multiple} onChange={pick} ref={inputRef} type="file" />
      <p className="pointer-events-none text-[11px]">{hint}</p>
    </div>
  );
}
