import Link from "next/link";

import { SiteFooter } from "@repo/ui/site-footer";
import { SiteHeader } from "@repo/ui/site-header";

import { Brand } from "../components/brand";
import Silk from "../components/Silk";
import { navLinks, tools } from "../lib/tools";

const workflow = ["Upload or drag your file into the tool.", "Adjust the options and click Run.", "Download the finished file."];

const cellClass =
  "relative flex min-h-[89px] w-full items-center justify-center border-b border-r border-line px-2 text-center text-[10px] text-tool transition-colors hover:border-white hover:bg-white hover:text-black focus-visible:border-white focus-visible:bg-white focus-visible:text-black focus-visible:outline-none sm:text-[11px]";

export default function Home() {
  return (
    <main className="flex min-h-screen w-full flex-col items-center bg-black font-mono text-[#d5d2d2]">
      <div className="w-full max-w-[684px] border-x border-line">
        <SiteHeader brand={<Brand />} links={navLinks} />

        <section className="relative flex min-h-[280px] flex-col items-center justify-center overflow-hidden border-b border-line px-5 py-16 text-start">
          <div aria-hidden className="pointer-events-none absolute inset-0 z-0  motion-reduce:hidden">
            <Silk speed={15} scale={1} color="#363636" noiseIntensity={1.0} rotation={0} />
          </div>

          <div className="relative z-10 w-full max-w-[500px]">
            <h1 className="text-[21px] font-bold leading-[1.35] tracking-[-0.045em] text-[#e7e3e3] sm:text-[24px]">The open source PDF &amp; Document Tools</h1>
            <p className=" text-[10px] text-[#969191] sm:text-[11px]">Gidocs brings useful tools for PDF files and documents into one simple place. Choose what you need and follow the steps.</p>
          </div>
        </section>

        <section className="border-b border-line px-5 py-5 sm:px-[50px] sm:py-[20px]" id="workflow">
          <div className="w-full max-w-[584px] text-left">
            <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">How it works</h2>
            <ul className="mt-2 space-y-0 text-[10px] leading-[1.8] text-[#a29d9d] sm:text-[11px]">
              {workflow.map((item, i) => (
                <li className="flex gap-2" key={item}>
                  <span className="shrink-0 font-bold  text-[#f3ab00]" aria-hidden="true">
                    {"[+]"}
                  </span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </div>
        </section>

        <section id="tools" aria-label="PDF tools">
          <ul className="grid grid-cols-2 border-l border-t border-line sm:grid-cols-5">
            {tools.map((tool) => {
              const badge = !tool.ready ? "soon" : tool.badge;
              const content = (
                <>
                  {tool.label}
                  {badge ?
                    <span className="absolute right-2 top-2 text-[9px] text-accent">[{badge}]</span>
                  : null}
                </>
              );
              return (
                <li key={tool.slug}>
                  {tool.ready ?
                    <Link className={cellClass} href={`/tools/${tool.slug}`}>
                      {content}
                    </Link>
                  : <span aria-disabled="true" className="relative flex min-h-[89px] w-full cursor-not-allowed items-center justify-center border-b border-r border-line px-2 text-center text-[10px] text-tool opacity-40 sm:text-[11px]">
                      {content}
                    </span>
                  }
                </li>
              );
            })}
          </ul>
        </section>

        <section className="min-h-[145px] border-t border-line px-5 py-8 sm:px-[50px] sm:py-5" id="privacy">
          <div className="mx-auto w-full max-w-[584px] text-left">
            <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">Privacy first</h2>
            <p className="mt-2 flex gap-2 text-[10px] leading-[1.9] text-[#a29d9d] sm:text-[11px]">
              <span className="shrink-0 font-bold text-[#f3ab00]" aria-hidden="true">
                [+]
              </span>
              <span>
                Gidocs is designed to help you work with sensitive documents. It has no database and keeps no files after your download. Read our{" "}
                <Link className="text-[#f3ab00] underline underline-offset-2 transition-colors hover:text-[#ffd36a]" href="/about">
                  project details
                </Link>{" "}
                to learn more.
              </span>
            </p>
          </div>
        </section>

        <SiteFooter />
      </div>
    </main>
  );
}
