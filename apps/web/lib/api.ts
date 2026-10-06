export interface RunResult {
  blob: Blob;
  filename: string;
}

export async function runTool(endpoint: string, files: File[], fields: Record<string, string>): Promise<RunResult> {
  const form = new FormData();
  for (const file of files) form.append("files", file);
  for (const [key, value] of Object.entries(fields)) form.append(key, value);

  const response = await fetch(endpoint, { method: "POST", body: form });
  if (!response.ok) {
    const raw = (await response.text()).trim();
    const message = raw && !raw.startsWith("<") ? raw : `Request failed (${response.status})`;
    throw new Error(message);
  }

  const blob = await response.blob();
  const filename = filenameFromDisposition(response.headers.get("Content-Disposition")) || `result-${Date.now()}.pdf`;
  return { blob, filename };
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
