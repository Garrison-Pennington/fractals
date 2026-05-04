// Pure-JS fallback implementation of the WASM API.
//
// The component prefers WASM (far faster and exactly matches the Rust core),
// but we ship this fallback so that:
//   1. During SSR / module eval time, `require`ing the component doesn't blow
//      up even though WASM isn't available.
//   2. If a host app hasn't set COOP/COEP headers or misconfigures the asset
//      path, the explorer still renders (just slower).
//
// Kept in sync with fractals-core semantics. When the Rust core changes, this
// file must be updated too — the wasm integration tests and the core unit
// tests are the source of truth.

import type { CustomFractalDef, FractalKind, GeneratedNote } from "./types.js";

function iterMandelbrot(cr: number, ci: number, maxIter: number) {
  let zr = 0,
    zi = 0;
  for (let i = 0; i < maxIter; i++) {
    const zr2 = zr * zr;
    const zi2 = zi * zi;
    if (zr2 + zi2 > 4) return { iter: i, mag: Math.sqrt(zr2 + zi2), escaped: true };
    const nzr = zr2 - zi2 + cr;
    const nzi = 2 * zr * zi + ci;
    zr = nzr;
    zi = nzi;
  }
  return { iter: maxIter, mag: Math.sqrt(zr * zr + zi * zi), escaped: false };
}

function iterJulia(zr0: number, zi0: number, cr: number, ci: number, maxIter: number) {
  let zr = zr0,
    zi = zi0;
  for (let i = 0; i < maxIter; i++) {
    const zr2 = zr * zr;
    const zi2 = zi * zi;
    if (zr2 + zi2 > 4) return { iter: i, mag: Math.sqrt(zr2 + zi2), escaped: true };
    const nzr = zr2 - zi2 + cr;
    const nzi = 2 * zr * zi + ci;
    zr = nzr;
    zi = nzi;
  }
  return { iter: maxIter, mag: Math.sqrt(zr * zr + zi * zi), escaped: false };
}

function iterBurningShip(cr: number, ci: number, maxIter: number) {
  let zr = 0,
    zi = 0;
  for (let i = 0; i < maxIter; i++) {
    const zr2 = zr * zr;
    const zi2 = zi * zi;
    if (zr2 + zi2 > 4) return { iter: i, mag: Math.sqrt(zr2 + zi2), escaped: true };
    const azr = Math.abs(zr);
    const azi = Math.abs(zi);
    const nzr = azr * azr - azi * azi + cr;
    const nzi = 2 * azr * azi + ci;
    zr = nzr;
    zi = nzi;
  }
  return { iter: maxIter, mag: Math.sqrt(zr * zr + zi * zi), escaped: false };
}

function iterNewton(cr: number, ci: number, maxIter: number) {
  // z^3 - 1 = 0 roots.
  const t = (Math.PI * 2) / 3;
  const roots = [
    [1, 0],
    [Math.cos(t), Math.sin(t)],
    [Math.cos(2 * t), Math.sin(2 * t)],
  ];
  let zr = cr,
    zi = ci;
  if (zr * zr + zi * zi < 1e-20) return { iter: maxIter, mag: 0, escaped: false };
  for (let i = 0; i < maxIter; i++) {
    // z^2
    const zr2 = zr * zr - zi * zi;
    const zi2 = 2 * zr * zi;
    // z^3
    const zr3 = zr2 * zr - zi2 * zi;
    const zi3 = zr2 * zi + zi2 * zr;
    // num = z^3 - 1
    const nr = zr3 - 1;
    const ni = zi3;
    // den = 3z^2
    const dr = 3 * zr2;
    const di = 3 * zi2;
    const denNorm = dr * dr + di * di;
    if (denNorm < 1e-24) return { iter: maxIter, mag: 0, escaped: false };
    const qr = (nr * dr + ni * di) / denNorm;
    const qi = (ni * dr - nr * di) / denNorm;
    zr -= qr;
    zi -= qi;
    for (let r = 0; r < 3; r++) {
      const dx = zr - roots[r][0];
      const dy = zi - roots[r][1];
      if (dx * dx + dy * dy < 1e-12) return { iter: i, mag: r, escaped: true };
    }
  }
  return { iter: maxIter, mag: 3, escaped: false };
}

function smooth(iter: number, mag: number, escaped: boolean, maxIter: number, escapeRadius = 2) {
  if (!escaped || iter >= maxIter) return maxIter;
  const logZn = Math.log(Math.max(1, mag));
  const nu = Math.log(logZn / Math.log(escapeRadius)) / Math.LN2;
  return Math.max(0, Math.min(maxIter, iter + 1 - nu));
}

// ---------------------------------------------------------------------------
// Custom fractal compiler
// ---------------------------------------------------------------------------
type IterFn = (cr: number, ci: number, re: number, im: number, maxIter: number) =>
  { iter: number; mag: number; escaped: boolean };

let _customCacheKey = "";
let _customCacheFn: IterFn | null = null;

export function compileCustomFractal(def: CustomFractalDef): IterFn {
  const key = [def.initZr, def.initZi, def.stepZr, def.stepZi, def.escapeExpr].join("\0");
  if (key === _customCacheKey && _customCacheFn) return _customCacheFn;
  // Build the function body.  Variables available to the user expressions:
  //   cr, ci  – c parameter (Julia) or pixel coordinate (Mandelbrot-like)
  //   re, im  – pixel coordinate in the complex plane
  //   zr, zi  – current iterate
  const body = `
    var zr = ${def.initZr}, zi = ${def.initZi};
    for (var i = 0; i < maxIter; i++) {
      if (${def.escapeExpr}) return { iter: i, mag: Math.sqrt(zr*zr + zi*zi), escaped: true };
      var nzr = ${def.stepZr};
      var nzi = ${def.stepZi};
      zr = nzr; zi = nzi;
      if (!isFinite(zr) || !isFinite(zi)) return { iter: i, mag: Infinity, escaped: true };
    }
    return { iter: maxIter, mag: Math.sqrt(zr*zr + zi*zi), escaped: false };
  `;
  const fn = new Function("cr", "ci", "re", "im", "maxIter", body) as IterFn;
  // Smoke-test to surface syntax errors immediately.
  fn(0, 0, 0, 0, 1);
  _customCacheKey = key;
  _customCacheFn = fn;
  return fn;
}

// Current custom definition, set by the render entry points.
let _activeCustomDef: CustomFractalDef | undefined;

export function setActiveCustomDef(def: CustomFractalDef | undefined) {
  _activeCustomDef = def;
}

function computeAt(
  kind: FractalKind,
  cRe: number,
  cIm: number,
  re: number,
  im: number,
  maxIter: number
) {
  switch (kind) {
    case "julia":
      return iterJulia(re, im, cRe, cIm, maxIter);
    case "burning_ship":
      return iterBurningShip(re, im, maxIter);
    case "newton":
      return iterNewton(re, im, maxIter);
    case "custom": {
      if (!_activeCustomDef) return iterMandelbrot(re, im, maxIter);
      const fn = compileCustomFractal(_activeCustomDef);
      return fn(re, im, re, im, maxIter);
    }
    case "mandelbrot":
    default:
      return iterMandelbrot(re, im, maxIter);
  }
}

export function renderFractalJs(
  kind: FractalKind,
  cRe: number,
  cIm: number,
  xMin: number,
  xMax: number,
  yMin: number,
  yMax: number,
  width: number,
  height: number,
  maxIter: number
): Float32Array {
  const out = new Float32Array(width * height * 2);
  for (let py = 0; py < height; py++) {
    const v = (py + 0.5) / height;
    const y = yMin + v * (yMax - yMin);
    for (let px = 0; px < width; px++) {
      const u = (px + 0.5) / width;
      const x = xMin + u * (xMax - xMin);
      const r = computeAt(kind, cRe, cIm, x, y, maxIter);
      const s = smooth(r.iter, r.mag, r.escaped, maxIter) / maxIter;
      const o = (py * width + px) * 2;
      out[o] = s;
      out[o + 1] = r.mag;
    }
  }
  return out;
}

export function renderProgressiveJs(
  kind: FractalKind,
  cRe: number,
  cIm: number,
  xMin: number,
  xMax: number,
  yMin: number,
  yMax: number,
  width: number,
  height: number,
  maxIter: number,
  stride: number
): Float32Array {
  const sw = Math.ceil(width / stride);
  const sh = Math.ceil(height / stride);
  const coarse = new Float32Array(sw * sh * 2);
  for (let sy = 0; sy < sh; sy++) {
    const py = Math.min(height - 1, sy * stride);
    const v = (py + 0.5) / height;
    const y = yMin + v * (yMax - yMin);
    for (let sx = 0; sx < sw; sx++) {
      const px = Math.min(width - 1, sx * stride);
      const u = (px + 0.5) / width;
      const x = xMin + u * (xMax - xMin);
      const r = computeAt(kind, cRe, cIm, x, y, maxIter);
      const s = smooth(r.iter, r.mag, r.escaped, maxIter) / maxIter;
      const o = (sy * sw + sx) * 2;
      coarse[o] = s;
      coarse[o + 1] = r.mag;
    }
  }
  const out = new Float32Array(width * height * 2);
  for (let py = 0; py < height; py++) {
    const sy = Math.min(sh - 1, Math.floor(py / stride));
    for (let px = 0; px < width; px++) {
      const sx = Math.min(sw - 1, Math.floor(px / stride));
      const src = (sy * sw + sx) * 2;
      const dst = (py * width + px) * 2;
      out[dst] = coarse[src];
      out[dst + 1] = coarse[src + 1];
    }
  }
  return out;
}

export function computePointJs(
  kind: FractalKind,
  cRe: number,
  cIm: number,
  re: number,
  im: number,
  maxIter: number
) {
  const r = computeAt(kind, cRe, cIm, re, im, maxIter);
  return {
    re,
    im,
    iterations: r.iter,
    escaped: r.escaped,
    smooth: smooth(r.iter, r.mag, r.escaped, maxIter) / maxIter,
  };
}

export function defaultBoundsJs(kind: FractalKind) {
  switch (kind) {
    case "julia":
      return { x_min: -1.5, x_max: 1.5, y_min: -1.5, y_max: 1.5 };
    case "burning_ship":
      return { x_min: -2.2, x_max: 1.5, y_min: -2.0, y_max: 1.2 };
    case "newton":
      return { x_min: -2.0, x_max: 2.0, y_min: -2.0, y_max: 2.0 };
    case "custom":
      if (_activeCustomDef) return { ...(_activeCustomDef.defaultBounds) };
      return { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 };
    case "mandelbrot":
    default:
      return { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 };
  }
}

// Music generation — mirrors the Rust sequence semantics closely enough for
// fallback playback. When WASM is available we call it instead.
const SCALE_INTERVALS: Record<string, number[]> = {
  major: [2, 4, 5, 7, 9, 11, 12],
  minor: [2, 3, 5, 7, 8, 10, 12],
  major_pentatonic: [2, 4, 7, 9, 12],
  minor_pentatonic: [3, 5, 7, 10, 12],
  pentatonic: [3, 5, 7, 10, 12],
  blues: [3, 5, 6, 7, 10, 12],
};

function scaleDegree(rootMidi: number, scale: string, degree: number): number {
  const intervals = SCALE_INTERVALS[scale] ?? SCALE_INTERVALS.major;
  const len = intervals.length;
  // degree is 0-based; 0 = root, 1 = first interval, etc.
  const octaves = Math.floor(degree / (len + 1));
  const idx = ((degree % (len + 1)) + (len + 1)) % (len + 1);
  const base = idx === 0 ? rootMidi : rootMidi + intervals[idx - 1];
  return base + octaves * 12;
}

const PC_NAMES = ["C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"];

function toneName(midi: number): string {
  const pc = ((midi - 12) % 12 + 12) % 12;
  const oct = Math.floor((midi - 12) / 12);
  return `${PC_NAMES[pc]}${oct}`;
}

export function generateMusicJs(cfg: {
  kind: FractalKind;
  c_re: number;
  c_im: number;
  scale: string;
  root_midi: number;
  sequence_type: string;
  length: number;
  max_iter: number;
  x_min: number;
  x_max: number;
  y_min: number;
  y_max: number;
}): GeneratedNote[] {
  const notes: GeneratedNote[] = [];
  const len = Math.max(1, cfg.length);
  if (cfg.sequence_type === "region") {
    const regions = Math.max(2, Math.min(8, len));
    const rw = (cfg.x_max - cfg.x_min) / regions;
    const side = 8;
    for (let r = 0; notes.length < len; r++) {
      const rIdx = r % regions;
      let total = 0;
      let escapes = 0;
      let count = 0;
      for (let sy = 0; sy < side; sy++) {
        const v = (sy + 0.5) / side;
        const y = cfg.y_min + v * (cfg.y_max - cfg.y_min);
        for (let sx = 0; sx < side; sx++) {
          const u = (sx + 0.5) / side;
          const x = cfg.x_min + rIdx * rw + u * rw;
          const s = computeAt(cfg.kind, cfg.c_re, cfg.c_im, x, y, cfg.max_iter);
          total += s.iter;
          if (s.escaped) escapes++;
          count++;
        }
      }
      const avg = total / count;
      const esc = escapes / count;
      const norm = Math.min(1, avg / cfg.max_iter);
      const intervalsLen = (SCALE_INTERVALS[cfg.scale] ?? SCALE_INTERVALS.major).length + 1;
      const degree = Math.round(norm * (intervalsLen - 1));
      const rootMidi = scaleDegree(cfg.root_midi, cfg.scale, degree);
      // Triad: root / third / fifth within scale.
      const chord = [0, 2, 4].map(d => scaleDegree(rootMidi, "major", d));
      for (const pitch of chord) {
        if (notes.length >= len) break;
        notes.push({
          pitch,
          duration: 0.5,
          velocity: Math.max(60, Math.min(110, Math.round(70 + esc * 40))),
          name: toneName(pitch),
        });
      }
    }
  } else {
    const yMid = (cfg.y_min + cfg.y_max) / 2;
    const scaleLen = (SCALE_INTERVALS[cfg.scale] ?? SCALE_INTERVALS.major).length + 1;
    const span = scaleLen * 3;
    for (let i = 0; i < len; i++) {
      const t = len <= 1 ? 0 : i / (len - 1);
      const x = cfg.x_min + t * (cfg.x_max - cfg.x_min);
      const s = computeAt(cfg.kind, cfg.c_re, cfg.c_im, x, yMid, cfg.max_iter);
      const degree = cfg.max_iter === 0 ? 0 : Math.floor((s.iter * span) / cfg.max_iter) % span;
      const pitch = scaleDegree(cfg.root_midi, cfg.scale, degree);
      const velocity = Math.max(30, Math.min(127, Math.round((s.iter * 127) / Math.max(1, cfg.max_iter))));
      notes.push({ pitch, duration: 0.5, velocity, name: toneName(pitch) });
    }
  }
  return notes;
}
