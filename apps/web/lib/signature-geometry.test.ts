import { describe, expect, test } from "bun:test";

import { clampPosition, cssToPt, displayScale, ptToCss } from "./signature-geometry";

describe("displayScale", () => {
  test("menggunakan lebar CSS yang tampil, bukan ukuran bitmap canvas", () => {
    // Halaman 612 pt dirender jadi bitmap 857 px (1.4 px/pt) tapi hanya
    // tampil 584 px karena CSS w-full di dalam shell 684 px.
    expect(displayScale(584, 612)).toBeCloseTo(0.954, 3);
  });

  test("jatuh ke 1 saat lebar CSS belum diukur", () => {
    expect(displayScale(0, 612)).toBe(1);
    expect(displayScale(584, 0)).toBe(1);
  });
});

describe("konversi titik PDF <-> piksel CSS", () => {
  test("bolak-balik tanpa kehilangan titik", () => {
    const scale = displayScale(584, 612);
    expect(ptToCss(72, scale)).toBeCloseTo(68.7, 1);
    expect(cssToPt(68.7, scale)).toBe(72);
  });

  test("cssToPt membulatkan ke titik bulat", () => {
    expect(cssToPt(10.4, 1)).toBe(10);
    expect(cssToPt(10.6, 1)).toBe(11);
  });
});

describe("clampPosition", () => {
  test("menahan tanda tangan agar tetap di dalam halaman", () => {
    const page = { width: 584, height: 758 };
    expect(clampPosition(-40, -40, 200, 60, page)).toEqual({ x: 0, y: 0 });
    // x maksimum 584-200=384; y maksimum 758-60-24 (margin)=674.
    expect(clampPosition(999, 999, 200, 60, page)).toEqual({ x: 384, y: 674 });
  });

  test("posisi yang sudah sah tidak berubah", () => {
    const page = { width: 584, height: 758 };
    expect(clampPosition(100, 200, 200, 60, page)).toEqual({ x: 100, y: 200 });
  });

  test("halaman lebih kecil dari tanda tangan dijepit ke 0", () => {
    expect(clampPosition(10, 10, 800, 900, { width: 584, height: 758 })).toEqual({ x: 0, y: 0 });
  });
});
