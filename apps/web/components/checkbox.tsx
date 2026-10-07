"use client";

import * as React from "react";
import { Checkbox as BaseCheckbox } from "@base-ui/react/checkbox";

interface CheckboxProps {
  checked: boolean;
  label: React.ReactNode;
  className?: string;
  onCheckedChange?: (checked: boolean) => void;
}

export default function Checkbox({ checked, label, className, onCheckedChange }: CheckboxProps) {
  return (
    <label className={className ?? "flex items-center gap-2 text-sm font-normal text-neutral-950 dark:text-white"}>
      <BaseCheckbox.Root
        checked={checked}
        className="flex size-4 shrink-0 items-center justify-center border rounded-none p-0 border-neutral-950 bg-white text-white dark:border-white dark:bg-neutral-950 dark:text-neutral-950 data-checked:bg-neutral-950 data-checked:text-white dark:data-checked:bg-white dark:data-checked:text-neutral-950 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-neutral-950 dark:focus-visible:outline-white"
        onCheckedChange={(value) => onCheckedChange?.(value)}
      >
        <BaseCheckbox.Indicator className="flex data-unchecked:hidden">
          <CheckIcon />
        </BaseCheckbox.Indicator>
      </BaseCheckbox.Root>
      {label}
    </label>
  );
}

function CheckIcon(props: React.ComponentProps<"svg">) {
  return (
    <svg
      width="16"
      height="16"
      viewBox="0 0 16 16"
      fill="none"
      stroke="currentColor"
      {...props}
      style={{ display: "block", ...props.style }}
    >
      <path d="m2.5 8.5 4 4 7-9" />
    </svg>
  );
}
