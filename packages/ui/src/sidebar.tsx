"use client";

import { type JSX } from "react";

export interface SidebarLink {
  label: string;
  href: string;
  external?: boolean;
}

interface SidebarProps {
  links: SidebarLink[];
  open: boolean;
  onClose: () => void;
}

export function Sidebar({ links, open, onClose }: SidebarProps): JSX.Element {
  return (
    <>
      <div
        className={`fixed inset-0 z-40 bg-black/70 transition-opacity duration-200 ${open ? "opacity-100" : "pointer-events-none opacity-0"}`}
        onClick={onClose}
      />
      <aside
        id="mobile-menu"
        inert={!open}
        aria-label="Menu"
        className={`fixed right-0 top-0 z-50 flex h-dvh w-full max-w-[320px] flex-col border-l border-[#252525] bg-black transition-transform duration-200 ${open ? "translate-x-0" : "translate-x-full"}`}
      >
        <div className="flex min-h-[51px] items-center justify-between border-b border-[#252525] px-5">
          <span className="flex items-center gap-1.5 text-[17px] font-bold tracking-[-0.08em] text-[#f3f1f1]">
            <span className="text-[20px] leading-none" aria-hidden="true">
              ✚
            </span>
            Gidocs
          </span>
          <button
            aria-label="Close menu"
            className="text-[11px] text-[#b8b4b4] transition-colors hover:text-white"
            onClick={onClose}
            type="button"
          >
            ✕
          </button>
        </div>
        <nav className="flex flex-col" aria-label="Mobile navigation">
          {links.map((link) => (
            <a
              className="border-b border-[#252525] px-5 py-4 text-[12px] text-[#b8b4b4] transition-colors hover:text-white"
              href={link.href}
              key={link.label}
              onClick={onClose}
              {...(link.external ? { target: "_blank", rel: "noreferrer" } : {})}
            >
              {link.label}
            </a>
          ))}
        </nav>
      </aside>
    </>
  );
}
