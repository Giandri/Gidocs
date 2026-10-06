"use client";

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
          disabled
            ? "cursor-not-allowed border-line text-tool opacity-40"
            : "border-accent text-accent hover:bg-white hover:text-black focus-visible:bg-white focus-visible:text-black"
        }`}
        disabled={disabled}
        onClick={onClick}
        type="button"
      >
        {loading ? "Processing [ ... ]" : "Run"}
      </button>
      {disabled && reason ? <p className="mt-2 text-[10px] text-muted">{reason}</p> : null}
    </div>
  );
}
