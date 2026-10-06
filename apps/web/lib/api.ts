export interface RunResult {
  blob: Blob;
  filename: string;
  /** Ukuran file sebelum diproses (untuk menampilkan sebelum/sesudah). */
  beforeBytes?: number;
}

export async function runTool(endpoint: string, files: File[], fields: Record<string, string>): Promise<RunResult> {
  const form = new FormData();
  for (const file of files) form.append("files", file);
  for (const [key, value] of Object.entries(fields)) form.append(key, value);

  const response = await fetch(endpoint, { method: "POST", body: form });
  if (!response.ok) throw new Error(await failMessage(response));

  const blob = await response.blob();
  const filename = filenameFromDisposition(response.headers.get("Content-Disposition")) || `result-${Date.now()}.pdf`;
  return { blob, filename };
}

const POLL_INTERVAL_MS = 800;
const JOB_TIMEOUT_MS = 120_000;

/** Menjalankan tool async: enqueue job, polling status, lalu unduh hasilnya. */
export async function runJob(slug: string, files: File[]): Promise<RunResult> {
  const form = new FormData();
  form.append("tool", slug);
  for (const file of files) form.append("files", file);

  const start = await fetch("/api/jobs", { method: "POST", body: form });
  if (!start.ok) throw new Error(await failMessage(start));
  const { id } = (await start.json()) as { id: string };

  const deadline = Date.now() + JOB_TIMEOUT_MS;
  while (Date.now() < deadline) {
    await new Promise((resolve) => setTimeout(resolve, POLL_INTERVAL_MS));
    const status = await fetch(`/api/jobs/${encodeURIComponent(id)}`);
    if (!status.ok) throw new Error(await failMessage(status));
    const info = (await status.json()) as { status: string; error?: string; download?: string };
    if (info.status === "error") throw new Error(info.error || "Conversion failed.");
    if (info.status !== "done") continue;

    const download = await fetch(info.download ?? `/api/jobs/${encodeURIComponent(id)}/download`);
    if (!download.ok) throw new Error(await failMessage(download));
    const blob = await download.blob();
    const filename =
      filenameFromDisposition(download.headers.get("Content-Disposition")) || `${slug}.pdf`;
    return { blob, filename };
  }
  throw new Error("Conversion timed out.");
}

async function failMessage(response: Response): Promise<string> {
  const raw = (await response.text()).trim();
  return raw && !raw.startsWith("<") ? raw : `Request failed (${response.status})`;
}

export function downloadBlob({ blob, filename }: RunResult): void {
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  link.click();
  URL.revokeObjectURL(url);
}

function filenameFromDisposition(header: string | null): string | null {
  if (!header) return null;
  const star = /filename\*=UTF-8''([^;]+)/i.exec(header);
  if (star?.[1]) return decodeURIComponent(star[1].trim());
  const plain = /filename="?([^";]+)"?/i.exec(header);
  return plain?.[1] ? plain[1].trim() : null;
}
