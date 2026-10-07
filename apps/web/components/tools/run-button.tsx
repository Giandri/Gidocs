"use client";

import LatticeLoader from "../LatticeLoader";

interface RunButtonProps {
  disabled: boolean;
  loading: boolean;
  reason?: string;
  onClick: () => void;
}

export function RunButton({ disabled, loading, reason, onClick }: RunButtonProps) {
  return (
    <div>
      <button
        className={`min-h-[44px] w-full border text-[11px] transition-colors focus-visible:outline-none ${
          disabled ? "cursor-not-allowed border-line text-tool opacity-40" : "border-accent text-accent hover:bg-white hover:text-black focus-visible:bg-white focus-visible:text-black"
        }`}
        disabled={disabled}
        onClick={onClick}
        type="button">
        {loading ?
          <span className="pointer-events-none inline-flex justify-center">
            <LatticeLoader
              status="working"
              label="Processing"
              doneLabel="Done in"
              errorLabel="Failed after"
              pattern="orbit"
              grid={3}
              shape="round"
              doneColor="#22c55e"
              errorColor="#ef4444"
              cellSize={3}
              gap={2}
              fontSize={14}
              step={90}
              idleOpacity={0.15}
              glow={false}
              glowColor=""
              showTimer
              color="#f5f5f5"
            />
          </span>
        : "Run"}
      </button>
      {disabled && reason ?
        <p className="mt-2 text-[10px] text-muted">{reason}</p>
      : null}
    </div>
  );
}
