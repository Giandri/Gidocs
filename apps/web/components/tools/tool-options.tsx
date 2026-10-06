"use client";

import { type OptionKind } from "../../lib/tools";

interface ToolOptionsProps {
  kind: OptionKind;
  fields: Record<string, string>;
  onChange: (name: string, value: string) => void;
}

const inputClass =
  "bg-black border border-line px-2 py-1.5 text-[11px] text-ink focus:border-white focus:outline-none disabled:opacity-40";
const labelClass = "text-[10px] text-muted";

export function ToolOptions({ kind, fields, onChange }: ToolOptionsProps) {
  if (kind === "split") {
    return (
      <div className="space-y-2">
        <Radio
          checked={fields.mode !== "ranges"}
          label="Every page"
          name="mode"
          onChange={() => onChange("mode", "every")}
        />
        <Radio
          checked={fields.mode === "ranges"}
          label="Page ranges"
          name="mode"
          onChange={() => onChange("mode", "ranges")}
        />
        {fields.mode === "ranges" ? (
          <label className="flex flex-col gap-1">
            <span className={labelClass}>Ranges (for example 1-3,5,8-)</span>
            <input
              className={inputClass}
              onChange={(event) => onChange("ranges", event.target.value)}
              placeholder="1-3,5,8-"
              value={fields.ranges}
            />
          </label>
        ) : null}
      </div>
    );
  }

  if (kind === "compress") {
    return (
      <div className="space-y-2">
        {(["light", "medium", "strong"] as const).map((level) => (
          <Radio
            checked={fields.level === level}
            key={level}
            label={level.charAt(0).toUpperCase() + level.slice(1)}
            name="level"
            onChange={() => onChange("level", level)}
          />
        ))}
      </div>
    );
  }

  if (kind === "watermark") {
    return (
      <div className="space-y-3">
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Watermark text</span>
          <input
            className={inputClass}
            onChange={(event) => onChange("text", event.target.value)}
            value={fields.text}
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Position</span>
          <select
            className={inputClass}
            onChange={(event) => onChange("position", event.target.value)}
            value={fields.position}
          >
            <option value="center">Center</option>
            <option value="diagonal">Diagonal</option>
            <option value="top-left">Top left</option>
            <option value="bottom-right">Bottom right</option>
          </select>
        </label>
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Opacity ({fields.opacity}%)</span>
          <input
            className="w-full accent-[#f3ab00]"
            max="100"
            min="10"
            onChange={(event) => onChange("opacity", event.target.value)}
            type="range"
            value={fields.opacity}
          />
        </label>
      </div>
    );
  }

  if (kind === "protect") {
    return (
      <div className="space-y-3">
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Password</span>
          <input
            autoComplete="new-password"
            className={inputClass}
            onChange={(event) => onChange("password", event.target.value)}
            type="password"
            value={fields.password}
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className={labelClass}>Confirm password</span>
          <input
            autoComplete="new-password"
            className={inputClass}
            onChange={(event) => onChange("confirm", event.target.value)}
            type="password"
            value={fields.confirm}
          />
        </label>
        <p className="flex gap-2 text-[10px] leading-[1.8] text-muted">
          <span className="shrink-0 font-bold text-accent" aria-hidden="true">
            [!]
          </span>
          <span>We cannot recover a lost password.</span>
        </p>
      </div>
    );
  }

  if (kind === "page-size") {
    return (
      <div className="space-y-2">
        <Radio
          checked={fields.size !== "fit"}
          label="A4"
          name="size"
          onChange={() => onChange("size", "a4")}
        />
        <Radio
          checked={fields.size === "fit"}
          label="Fit image"
          name="size"
          onChange={() => onChange("size", "fit")}
        />
      </div>
    );
  }

  return null;
}

function Radio({
  checked,
  label,
  name,
  onChange,
}: {
  checked: boolean;
  label: string;
  name: string;
  onChange: () => void;
}) {
  return (
    <label className="flex cursor-pointer items-center gap-2 text-[11px] text-ink">
      <input checked={checked} className="accent-[#f3ab00]" name={name} onChange={onChange} type="radio" />
      {label}
    </label>
  );
}
