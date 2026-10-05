import { MobileNav } from "@repo/ui/mobile-nav";

const tools = ["DOCX to PDF", "Excel to PDF", "PPTX to PDF", "JPEG to PDF", "Compress PDF", "Split PDF", "Merge PDF", "Watermark PDF", "Protect PDF", "Sign PDF"];

const badges: Record<string, string> = { "DOCX to PDF": "most", "Protect PDF": "new", "Sign PDF": "new" };

const workflow = ["Choose the document or file you want to work with.", "Find a focused tool for your document task.", "Download your finished PDF, ready to use or share."];

const repository = "https://github.com/Giandri/Gidocs";

const navLinks = [
  { label: "GitHub", href: repository, external: true },
  { label: "Docs", href: "#workflow" },
  { label: "Stars", href: "#privacy" },
];

export default function Home() {
  return (
    <main className="flex min-h-screen w-full flex-col  items-center bg-black px- font-mono text-[#d5d2d2]">
      <div className="w-full max-w-[684px] border-x border-[#252525]">
        <header className="border-b border-[#252525]">
          <nav className="flex min-h-[51px] items-center justify-between px-5 sm:px-[30px]" aria-label="Main navigation">
            <a className="flex items-center gap-1.5 text-[17px] font-bold tracking-[-0.08em] text-[#f3f1f1]" href="#" aria-label="Gidocs home">
              <span className="text-[20px] leading-none" aria-hidden="true">
                ✚
              </span>
              Gidocs
            </a>
            <div className="hidden items-center gap-3 text-[10px] text-[#b8b4b4] sm:gap-6 sm:text-[11px] sm:flex">
              {navLinks.map((link) => (
                <a className="transition-colors hover:text-white" href={link.href} key={link.label} {...(link.external ? { target: "_blank", rel: "noreferrer" } : {})}>
                  {link.label}
                </a>
              ))}
            </div>
            <MobileNav links={navLinks} />
          </nav>
        </header>

        <section className="flex min-h-[280px] flex-col items-center justify-center border-b border-[#252525]  px-5 py-16 text-start">
          <div className="w-full max-w-[500px]">
            <h1 className="text-[21px] font-bold leading-[1.35] tracking-[-0.045em] text-[#e7e3e3] sm:text-[24px]">The open source PDF &amp; Document Tools</h1>
            <p className=" text-[10px] text-[#969191] sm:text-[11px]">Gidoc brings useful tools for PDF files and documents into one simple place. Choose what you need and follow the steps.</p>
          </div>
        </section>

        <section className="border-b border-[#252525] px-5 py-5 sm:px-[50px] sm:py-[20px]" id="workflow">
          <div className="w-full max-w-[584px] text-left p-10">
            <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">The open source AI coding agent</h2>
            <ul className="mt-2 space-y-1 text-[10px] leading-[1.8] text-[#a29d9d] sm:text-[11px]">
              {workflow.map((item) => (
                <li className="flex gap-2" key={item}>
                  <span className="shrink-0 font-bold  text-[#f3ab00]" aria-hidden="true">
                    [*]
                  </span>
                  <span>{item}</span>
                </li>
              ))}
            </ul>
          </div>
        </section>

        <section id="tools" aria-label="PDF tools">
          <ul className="grid grid-cols-2 border-l border-t border-[#252525] sm:grid-cols-5">
            {tools.map((tool) => (
              <li key={tool}>
                <button
                  className="relative flex min-h-[89px] w-full items-center justify-center border-b border-r border-[#252525] px-2 text-center text-[10px] text-[#c7c2c2] transition-colors hover:border-white hover:bg-white hover:text-black focus-visible:border-white focus-visible:bg-white focus-visible:text-black focus-visible:outline-none sm:text-[11px]"
                  type="button">
                  {tool}
                  {badges[tool] && <span className="absolute right-2 top-2 text-[9px] text-[#f3ab00]">[{badges[tool]}]</span>}
                </button>
              </li>
            ))}
          </ul>
        </section>

        <section className="min-h-[145px] border-t border-[#252525] px-5 py-8 sm:px-[50px] sm:py-10" id="privacy">
          <div className="mx-auto w-full max-w-[584px] text-left p-10">
            <h2 className="text-[11px] font-bold text-[#e4e0e0] sm:text-xs">Privacy first</h2>
            <p className="mt-3 flex gap-2 text-[10px] leading-[1.9] text-[#a29d9d] sm:text-[11px]">
              <span className="shrink-0 font-bold text-[#f3ab00]" aria-hidden="true">
                [+]
              </span>
              <span>
                Gidoc is designed to help you work with sensitive documents. Review our{" "}
                <a className="text-[#f3ab00] underline underline-offset-2 transition-colors hover:text-[#ffd36a]" href={`${repository}/blob/main/README.md`} target="_blank" rel="noreferrer">
                  project details
                </a>{" "}
                to learn more.
              </span>
            </p>
          </div>
        </section>

        <footer className="mt-10 flex min-h-[17px] w-full items-center justify-center border-t border-[#252525] px-4 text-[9px] text-[#797474]">©2026 Giandri</footer>
      </div>
    </main>
  );
}
