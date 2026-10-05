"use client";

import { type JSX, useEffect, useState } from "react";

import { type SidebarLink, Sidebar } from "./sidebar";

export function MobileNav({ links }: { links: SidebarLink[] }): JSX.Element {
  const [open, setOpen] = useState(false);

  useEffect(() => {
    if (!open) return;

    const previous = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };
    window.addEventListener("keydown", onKeyDown);

    return () => {
      document.body.style.overflow = previous;
      window.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  return (
    <>
      <button
        aria-controls="mobile-menu"
        aria-expanded={open}
        aria-label={open ? "Close menu" : "Open menu"}
        className="px-2 text-[11px] text-[#b8b4b4] transition-colors hover:text-white sm:hidden"
        onClick={() => setOpen((value) => !value)}
        type="button"
      >
        {open ? "✕" : "☰"}
      </button>
      <Sidebar links={links} onClose={() => setOpen(false)} open={open} />
    </>
  );
}
