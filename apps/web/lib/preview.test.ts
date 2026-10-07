import { describe, expect, test } from "bun:test";

import { canPreviewResult } from "./preview";

describe("canPreviewResult", () => {
  test("PDF bisa dipratinjau baik dari nama file maupun tipe", () => {
    expect(canPreviewResult("signed.pdf", "application/pdf")).toBe(true);
    expect(canPreviewResult("hasil", "application/pdf")).toBe(true);
    expect(canPreviewResult("SIGNED.PDF", "application/octet-stream")).toBe(true);
  });

  test("ZIP dan berkas lain tidak bisa dipratinjau", () => {
    expect(canPreviewResult("parts.zip", "application/zip")).toBe(false);
    expect(canPreviewResult("images.zip", "application/octet-stream")).toBe(false);
    expect(canPreviewResult("gambar.png", "image/png")).toBe(false);
    expect(canPreviewResult("", "")).toBe(false);
  });
});
