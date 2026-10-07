import { type ReactNode } from "react";

import { SiteFooter } from "@repo/ui/site-footer";
import { SiteHeader } from "@repo/ui/site-header";

import { navLinks, type Tool } from "../../lib/tools";
import { Brand } from "../brand";

export function ToolShell({ tool, children }: { tool: Tool; children: ReactNode }) {
  return (
    <main className="flex min-h-screen w-full flex-col items-center bg-black font-mono text-[#d5d2d2]">
      <div className="flex min-h-screen w-full max-w-[684px] flex-col border-x border-line">
        <SiteHeader brand={<Brand />} links={navLinks} />
        <section className="border-b border-line px-5 py-3 sm:px-[50px]">
          <a className="text-[10px] text-nav transition-colors hover:text-white sm:text-[11px]" href="/">
            {"<"} All tools
          </a>
        </section>
        <section className="border-b border-line px-5 py-6 sm:px-[50px]">
          <h1 className="text-[21px] font-bold leading-[1.35] tracking-[-0.045em] text-[#e7e3e3] sm:text-[24px]">
            {tool.title}
          </h1>
          <p className="mt-1 text-[10px] text-muted sm:text-[11px]">{tool.description}</p>
        </section>
        {children}
        <section className="border-b border-line px-5 py-3 sm:px-[50px]">
          <p className="flex gap-2 text-[10px] leading-[1.8] text-muted sm:text-[11px]">
            <span className="shrink-0 font-bold text-accent" aria-hidden="true">
              [+]
            </span>
            <span>Your files are processed temporarily and deleted right after the job finishes.</span>
          </p>
        </section>
        <div className="mt-auto">
          <SiteFooter />
        </div>
      </div>
    </main>
  );
}
