import { type ReactNode } from "react";

import { MobileNav } from "./mobile-nav";
import { type SidebarLink } from "./sidebar";

const updates = [
  { tag: "[new update]", text: "OCR PDF — recognize text in scanned PDFs" },
  { tag: "[new update]", text: "Sign PDF — draw or type a signature and place it on the page" },
  { tag: "[new update]", text: "PDF to Images — export every page as a PNG" },
  { tag: "[new update]", text: "Result preview — check the file before downloading it" },
  { tag: "[new update]", text: "Office to PDF — DOCX, Excel, and PPTX" },
  { tag: "[privacy]", text: "no database, no accounts, files deleted after download" },
];

export function SiteHeader({ brand, links }: { brand: ReactNode; links: SidebarLink[] }) {
  return (
    <>
      {/* sticky dengan latar solid: konten di bawahnya tidak boleh terlihat menembus header. */}
      <header className="sticky top-0 z-30 border-b border-[#252525] bg-black">
        <nav className="flex min-h-[51px] items-center justify-between border-b border-[#252525] px-5 sm:px-[30px]" aria-label="Main navigation">
          <a aria-label="Gidocs home" className="flex items-center gap-1.5 text-[17px] font-bold tracking-[-0.08em] text-[#f3f1f1]" href="/">
            {brand}
          </a>
          <div className="hidden items-center gap-3 text-[10px] text-[#b8b4b4] sm:gap-6 sm:text-[11px] sm:flex">
            {links.map((link) => (
              <a className="transition-colors hover:text-white" href={link.href} key={link.label} {...(link.external ? { target: "_blank", rel: "noreferrer" } : {})}>
                {link.label}
              </a>
            ))}
          </div>
          <MobileNav links={links} />
        </nav>
        {/*
          Pita informasi berjalan. Animasi murni CSS: konten digandakan lalu
          digeser -50% lewat keyframes `gidocs-marquee`. Salinan kedua
          disembunyikan dari pembaca layar dan animasi dimatikan saat pengguna
          meminta reduced motion. Garis pemisah navbar ada di border-b <nav>.
        */}
        <div className="overflow-hidden bg-black">
          <div className="flex w-max animate-[gidocs-marquee_38s_linear_infinite] motion-reduce:animate-none">
            {[false, true].map((duplicate) => (
              <p aria-hidden={duplicate || undefined} className="flex w-max shrink-0 items-center gap-5 px-3.5 py-2 text-[10px] text-muted sm:text-[10px]" key={duplicate ? "copy" : "original"}>
                {updates.map((item) => (
                  <span className="flex items-center gap-2 whitespace-nowrap" key={item.text}>
                    <span className="font-bold text-[#f3ab00]">{item.tag}</span>
                    <span>{item.text}</span>
                  </span>
                ))}
              </p>
            ))}
          </div>
        </div>
      </header>
    </>
  );
}
