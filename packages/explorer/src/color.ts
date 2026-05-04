// OKLAB colour-space palette interpolation.
//
// Smooth-iteration counts from the WASM renderer are in [0, 1]. Each palette
// is defined as a list of sRGB stops; we convert them to OKLAB up front, then
// for every pixel we:
//   1. Find the pair of stops surrounding the input t.
//   2. Interpolate linearly in OKLAB.
//   3. Convert the result back to sRGB for the canvas.
//
// OKLAB interpolation is perceptually uniform: smooth gradients stay free of
// visible bands that linear sRGB interpolation can introduce.

import type { PaletteName } from "./types.js";

type RGB = [number, number, number];
type OKLAB = [number, number, number];

const PALETTES: Record<PaletteName, RGB[]> = {
  // Classic blue-to-black Mandelbrot palette.
  classic: [
    [0, 0, 0],
    [32, 48, 120],
    [60, 120, 200],
    [200, 230, 255],
    [255, 255, 220],
    [180, 100, 40],
    [30, 20, 10],
    [0, 0, 0],
  ],
  // Fire: deep red -> orange -> white -> black for the interior.
  fire: [
    [0, 0, 0],
    [60, 0, 0],
    [180, 20, 0],
    [255, 120, 10],
    [255, 220, 80],
    [255, 255, 220],
    [100, 20, 0],
    [0, 0, 0],
  ],
  // Psychedelic rainbow.
  rainbow: [
    [0, 0, 0],
    [120, 0, 180],
    [0, 60, 220],
    [0, 200, 200],
    [0, 230, 60],
    [240, 230, 0],
    [240, 90, 0],
    [140, 0, 40],
    [0, 0, 0],
  ],
};

// sRGB <-> linear helpers.
function srgbToLinear(c: number): number {
  const v = c / 255;
  return v <= 0.04045 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4);
}
function linearToSrgb(c: number): number {
  const v = c <= 0.0031308 ? 12.92 * c : 1.055 * Math.pow(c, 1 / 2.4) - 0.055;
  return Math.max(0, Math.min(255, Math.round(v * 255)));
}

// Linear-sRGB to OKLab. Formulas from Björn Ottosson (2020).
function rgbToOklab(r: number, g: number, b: number): OKLAB {
  const lr = srgbToLinear(r);
  const lg = srgbToLinear(g);
  const lb = srgbToLinear(b);
  const l = 0.4122214708 * lr + 0.5363325363 * lg + 0.0514459929 * lb;
  const m = 0.2119034982 * lr + 0.6806995451 * lg + 0.1073969566 * lb;
  const s = 0.0883024619 * lr + 0.2817188376 * lg + 0.6299787005 * lb;
  const l_ = Math.cbrt(l);
  const m_ = Math.cbrt(m);
  const s_ = Math.cbrt(s);
  return [
    0.2104542553 * l_ + 0.793617785 * m_ - 0.0040720468 * s_,
    1.9779984951 * l_ - 2.428592205 * m_ + 0.4505937099 * s_,
    0.0259040371 * l_ + 0.7827717662 * m_ - 0.808675766 * s_,
  ];
}

function oklabToRgb(L: number, a: number, b: number): RGB {
  const l_ = L + 0.3963377774 * a + 0.2158037573 * b;
  const m_ = L - 0.1055613458 * a - 0.0638541728 * b;
  const s_ = L - 0.0894841775 * a - 1.291485548 * b;
  const l = l_ * l_ * l_;
  const m = m_ * m_ * m_;
  const s = s_ * s_ * s_;
  const r = 4.0767416621 * l - 3.3077115913 * m + 0.2309699292 * s;
  const g = -1.2684380046 * l + 2.6097574011 * m - 0.3413193965 * s;
  const bb = -0.0041960863 * l - 0.7034186147 * m + 1.707614701 * s;
  return [linearToSrgb(r), linearToSrgb(g), linearToSrgb(bb)];
}

// Pre-baked lookup tables (one 4096-entry LUT per palette). For each palette
// we precompute the sRGB value at 4096 evenly spaced t values so that the hot
// pixel loop becomes a table lookup + bilinear interpolation between rows.
const LUT_SIZE = 4096;
const LUTS = new Map<PaletteName, Uint8ClampedArray>();

function buildLut(name: PaletteName): Uint8ClampedArray {
  const stops = PALETTES[name];
  const labStops: OKLAB[] = stops.map(([r, g, b]) => rgbToOklab(r, g, b));
  const out = new Uint8ClampedArray(LUT_SIZE * 3);
  const segments = stops.length - 1;
  for (let i = 0; i < LUT_SIZE; i++) {
    const t = i / (LUT_SIZE - 1);
    const x = t * segments;
    const idx = Math.min(segments - 1, Math.floor(x));
    const frac = x - idx;
    const a = labStops[idx];
    const b = labStops[idx + 1];
    const L = a[0] + (b[0] - a[0]) * frac;
    const aa = a[1] + (b[1] - a[1]) * frac;
    const bb = a[2] + (b[2] - a[2]) * frac;
    const [r, g, bl] = oklabToRgb(L, aa, bb);
    out[i * 3] = r;
    out[i * 3 + 1] = g;
    out[i * 3 + 2] = bl;
  }
  return out;
}

export function getLut(name: PaletteName): Uint8ClampedArray {
  let lut = LUTS.get(name);
  if (!lut) {
    lut = buildLut(name);
    LUTS.set(name, lut);
  }
  return lut;
}

export const PALETTE_NAMES: PaletteName[] = ["classic", "fire", "rainbow"];

/**
 * Colorize a smooth-iteration buffer into an ImageData-ready Uint8ClampedArray.
 *
 * - `smoothBuf` layout: [t0, extra0, t1, extra1, ...] where t is in [0,1] and
 *   extra encodes extra per-pixel data (e.g. Newton basin id).
 * - Pixels that never escaped (t >= ~1) are rendered black (inside the set).
 * - For Newton (`isNewton`), extra is treated as a basin index 0..2 and we
 *   shift the hue by basin so each root paints a distinct colour family.
 */
export function colorize(
  smoothBuf: Float32Array,
  width: number,
  height: number,
  palette: PaletteName,
  isNewton: boolean
): Uint8ClampedArray {
  const lut = getLut(palette);
  const out = new Uint8ClampedArray(width * height * 4);
  const n = width * height;
  for (let i = 0; i < n; i++) {
    const t = smoothBuf[i * 2];
    const extra = smoothBuf[i * 2 + 1];
    const o = i * 4;
    // Inside the set / non-converged: black.
    if (t >= 0.999) {
      out[o] = 0;
      out[o + 1] = 0;
      out[o + 2] = 0;
      out[o + 3] = 255;
      continue;
    }
    let lutT = t;
    if (isNewton) {
      // Stack basins onto different portions of the palette.
      const basin = Math.min(2, Math.max(0, Math.round(extra)));
      lutT = basin / 3 + t / 3;
    } else {
      // Accent: boost contrast via a gentle gamma curve.
      lutT = Math.pow(t, 0.7);
    }
    const idx = Math.max(0, Math.min(LUT_SIZE - 1, Math.floor(lutT * (LUT_SIZE - 1))));
    const lo = idx * 3;
    out[o] = lut[lo];
    out[o + 1] = lut[lo + 1];
    out[o + 2] = lut[lo + 2];
    out[o + 3] = 255;
  }
  return out;
}
