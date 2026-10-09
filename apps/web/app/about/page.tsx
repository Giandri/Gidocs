import type { Metadata } from "next";

import { SiteFooter } from "@repo/ui/site-footer";
import { SiteHeader } from "@repo/ui/site-header";

import { Brand } from "../../components/brand";
import { navLinks, repository } from "../../lib/tools";
import Silk from "@/components/Silk";

export const metadata: Metadata = {
  title: "About - Gidocs",
  description: "Why Gidocs exists: open source PDF tools that keep sensitive documents out of public storage by using no database and no permanent file storage.",
};

const sections = [
  {
    heading: "Why this project exists",
    body: [
      "Many PDF tools ask you to upload a contract, a payslip, a national ID, a medical record, or a signed contract to somebody else's server, then keep it. A file that lands in a public bucket can leak through a bug, a misconfigured setting, or an employee with a bad day.",
      "Gidocs takes the opposite position. The project is built so that your document is the least interesting thing about it: it is processed, handed back to you, and forgotten.",
    ],
  },
  {
    heading: "How your files are protected",
    body: [
      "No database. Gidocs has no database and no permanent file storage. Nothing you upload becomes a record that can be browsed later.",
      "No accounts. There is no sign-up, no login, no history, and no file library. Without an account, there is nothing for a stranger to browse and nothing to leak.",
      "Temporary processing only. Every request gets its own temporary working folder. When the response is sent, that folder is deleted. Background jobs (Office conversion, PDF to images, OCR) expire on a timer as well, and a cleaner removes anything left behind.",
      "Files are treated as hostile input. Uploads are accepted based on their real file signature, not the name or the browser supplied type. Your original file name is never used as a path, and neither file contents nor original names are written to the logs.",
      "No third-party tracking. There are no analytics scripts and no trackers following you across the web.",
      "Preview stays on your device. When a tool needs to show you a page, for example the Sign PDF preview, that rendering happens in your own browser. The preview is never uploaded.",
    ],
  },
  {
    heading: "What Gidocs deliberately does not do",
    body: ["No accounts and no billing. No saved history. No shareable links. No long-term storage of your documents. These are not features waiting to be unlocked later; they are the reason the project is built this way."],
  },
  {
    heading: "Honest limits",
    body: [
      "To convert a file, the file does pass through the server that runs the tool. That is unavoidable for PDF conversion: the work happens on a machine. What we guarantee is how short that stay is and what is left behind afterwards.",
      "If your documents are sensitive enough that even a brief server visit is not acceptable, self-host Gidocs and run it where you control the machine. The whole project is open source precisely so you can do that.",
    ],
  },
];

export default function AboutPage() {
  return (
    <main className="flex min-h-screen w-full flex-col items-center bg-black font-mono text-[#d5d2d2]">
      <div className="w-full max-w-[684px] border-x border-line">
        <SiteHeader brand={<Brand />} links={navLinks} />
        <section className="border-b border-line px-5 py-3 sm:px-[50px]">
          <a className="text-[10px] text-nav transition-colors hover:text-white sm:text-[11px]" href="/">
            {"<"} Back to tools
          </a>
        </section>
        <section className="border-b border-line px-5 py-12 text-start sm:px-[50px] sm:py-14">
          <div className="w-full max-w-[584px]">
            <h1 className="text-[21px] font-bold leading-[1.35] tracking-[-0.045em] text-[#e7e3e3] sm:text-[24px]">About Gidocs</h1>
            <p className="mt-2 text-[10px] leading-[1.9] text-[#969191] sm:text-[11px]">Open source PDF and document tools that keep sensitive documents out of public storage.</p>
          </div>
        </section>

        {sections.map((section) => (
          <section className="border-b border-line px-5 py-6 sm:px-[50px] sm:py-5" key={section.heading}>
            <div className="w-full max-w-[584px] text-left">
              <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">{section.heading}</h2>
              {section.body.map((paragraph) => (
                <p className="mt-2 flex gap-2 text-[10px] leading-[1.9] text-[#a29d9d] sm:text-[11px]" key={paragraph}>
                  <span className="shrink-0 font-bold text-justify text-[#f3ab00]" aria-hidden="true">
                    [+]
                  </span>
                  <span className="text-justify">{paragraph}</span>
                </p>
              ))}
            </div>
          </section>
        ))}

        <section className="border-b border-line px-5 py-6 sm:px-[50px] sm:py-5">
          <div className="w-full max-w-[584px] text-left">
            <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">Open source</h2>
            <p className="mt-2 flex gap-2 text-[10px] text-justify leading-[1.9] text-[#a29d9d] sm:text-[11px]">
              <span className="shrink-0 font-bold text-[#f3ab00]" aria-hidden="true">
                [+]
              </span>
              <span>
                Read the source, run it yourself, or help improve it on{" "}
                <a className="text-[#f3ab00] underline underline-offset-2 transition-colors hover:text-[#ffd36a]" href={repository} rel="noreferrer" target="_blank">
                  GitHub
                </a>
                . Gidocs is released under the AGPL-3.0 license.
              </span>
            </p>
          </div>
        </section>

        <SiteFooter />
      </div>
    </main>
  );
}
