import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { ToolPage } from "../../../components/tools/tool-page";
import { ToolShell } from "../../../components/tools/tool-shell";
import { getTool, tools } from "../../../lib/tools";

interface PageProps {
  params: Promise<{ slug: string }>;
}

export function generateStaticParams() {
  return tools.map((tool) => ({ slug: tool.slug }));
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { slug } = await params;
  const tool = getTool(slug);
  if (!tool) return {};
  return {
    title: `${tool.title} — Gidocs`,
    description: tool.description,
  };
}

export default async function ToolRoute({ params }: PageProps) {
  const { slug } = await params;
  const tool = getTool(slug);
  if (!tool) notFound();
  return (
    <ToolShell tool={tool}>
      <ToolPage tool={tool} />
    </ToolShell>
  );
}
