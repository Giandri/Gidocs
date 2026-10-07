/** Apakah hasil tool bisa ditampilkan di halaman tanpa diunduh. */
export function canPreviewResult(filename: string, mime: string): boolean {
  return filename.toLowerCase().endsWith(".pdf") || mime === "application/pdf";
}
