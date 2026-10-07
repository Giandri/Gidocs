export type OptionKind = "none" | "split" | "compress" | "watermark" | "protect" | "page-size";

export interface Tool {
  slug: string;
  label: string;
  title: string;
  description: string;
  accept: string;
  multiple: boolean;
  minFiles: number;
  maxFiles: number;
  endpoint: string;
  ready: boolean;
  /** true = diproses lewat antrian jobs (async), bukan request sinkron. */
  async?: boolean;
  badge?: string;
  options: OptionKind;
  /** Ekstensi file unduhan (default ".pdf"). */
  outputExt?: string;
}

export const repository = "https://github.com/Giandri/Gidocs";

export const navLinks = [
  { label: "GitHub", href: repository, external: true },
  { label: "Docs", href: "/#workflow" },
  { label: "Stars", href: `${repository}/stargazers`, external: true },
];

export const tools: Tool[] = [
  {
    slug: "docx-to-pdf",
    label: "DOCX to PDF",
    title: "DOCX to PDF",
    description: "Convert a Word document to PDF.",
    accept: ".docx",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/jobs",
    ready: true,
    async: true,
    badge: "most",
    options: "none",
  },
  {
    slug: "excel-to-pdf",
    label: "Excel to PDF",
    title: "Excel to PDF",
    description: "Convert a spreadsheet to PDF.",
    accept: ".xlsx",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/jobs",
    ready: true,
    async: true,
    options: "none",
  },
  {
    slug: "pptx-to-pdf",
    label: "PPTX to PDF",
    title: "PPTX to PDF",
    description: "Convert a presentation to PDF.",
    accept: ".pptx",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/jobs",
    ready: true,
    async: true,
    options: "none",
  },
  {
    slug: "jpeg-to-pdf",
    label: "JPEG to PDF",
    title: "JPEG to PDF",
    description: "Turn JPG or PNG images into one PDF.",
    accept: ".jpg,.jpeg,.png",
    multiple: true,
    minFiles: 1,
    maxFiles: 20,
    endpoint: "/api/images-to-pdf",
    ready: true,
    options: "page-size",
  },
  {
    slug: "compress-pdf",
    label: "Compress PDF",
    title: "Compress PDF",
    description: "Make a PDF smaller.",
    accept: ".pdf",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/compress",
    ready: true,
    options: "compress",
  },
  {
    slug: "split-pdf",
    label: "Split PDF",
    title: "Split PDF",
    description: "Split a PDF into separate files.",
    accept: ".pdf",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/split",
    ready: true,
    options: "split",
    outputExt: ".zip",
  },
  {
    slug: "merge-pdf",
    label: "Merge PDF",
    title: "Merge PDF",
    description: "Combine multiple PDF files into one.",
    accept: ".pdf",
    multiple: true,
    minFiles: 2,
    maxFiles: 20,
    endpoint: "/api/merge",
    ready: true,
    badge: "new",
    options: "none",
  },
  {
    slug: "watermark-pdf",
    label: "Watermark PDF",
    title: "Watermark PDF",
    description: "Stamp text across every page.",
    accept: ".pdf",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/watermark",
    ready: true,
    options: "watermark",
  },
  {
    slug: "protect-pdf",
    label: "Protect PDF",
    title: "Protect PDF",
    description: "Encrypt a PDF with a password.",
    accept: ".pdf",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "/api/protect",
    ready: true,
    options: "protect",
  },
  {
    slug: "sign-pdf",
    label: "Sign PDF",
    title: "Sign PDF",
    description: "Add a signature to a PDF.",
    accept: ".pdf",
    multiple: false,
    minFiles: 1,
    maxFiles: 1,
    endpoint: "",
    ready: false,
    badge: "new",
    options: "none",
  },
];

export function getTool(slug: string): Tool | undefined {
  return tools.find((tool) => tool.slug === slug);
}
