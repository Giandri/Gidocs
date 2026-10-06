import { type ReactNode } from "react";

import { MobileNav } from "./mobile-nav";
import { type SidebarLink } from "./sidebar";

export function SiteHeader({ brand, links }: { brand: ReactNode; links: SidebarLink[] }) {
  return (
    <header className="border-b border-[#252525]">
      <nav className="flex min-h-[51px] items-center justify-between px-5 sm:px-[30px]" aria-label="Main navigation">
        <a
          aria-label="Gidocs home"
          className="flex items-center gap-1.5 text-[17px] font-bold tracking-[-0.08em] text-[#f3f1f1]"
          href="/"
        >
          {brand}
        </a>
        <div className="hidden items-center gap-3 text-[10px] text-[#b8b4b4] sm:gap-6 sm:text-[11px] sm:flex">
          {links.map((link) => (
            <a
              className="transition-colors hover:text-white"
              href={link.href}
              key={link.label}
              {...(link.external ? { target: "_blank", rel: "noreferrer" } : {})}
            >
              {link.label}
            </a>
          ))}
        </div>
        <MobileNav links={links} />
      </nav>
    </header>
  );
}
