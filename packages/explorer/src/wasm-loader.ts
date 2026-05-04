// WASM loader that tries the Rust module first and gracefully falls back to
// the JS implementation when running in SSR / when the .wasm can't be served.
//
// The consumer can override `assetPath` via a prop on <FractalExplorer /> to
// point at a custom location (e.g. a CDN). Default is `./wasm/` relative to
// the importing module.

import type { FractalKind, GeneratedNote, PointInfo } from "./types.js";
import {
  computePointJs,
  defaultBoundsJs,
  generateMusicJs,
  renderFractalJs,
  renderProgressiveJs,
} from "./fallback.js";

export { setActiveCustomDef } from "./fallback.js";

export interface FractalsApi {
  kind: "wasm" | "js";
  render(
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
  ): Float32Array;
  renderProgressive(
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
  ): Float32Array;
  computePoint(
    kind: FractalKind,
    cRe: number,
    cIm: number,
    re: number,
    im: number,
    maxIter: number
  ): PointInfo;
  defaultBounds(kind: FractalKind): { x_min: number; x_max: number; y_min: number; y_max: number };
  generateMusic(cfg: {
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
  }): GeneratedNote[];
}

// JS fallback always available — used immediately while WASM is initializing,
// during SSR, or if WASM load fails.
export const jsApi: FractalsApi = {
  kind: "js",
  render: renderFractalJs,
  renderProgressive: renderProgressiveJs,
  computePoint: (kind, cRe, cIm, re, im, maxIter) =>
    computePointJs(kind, cRe, cIm, re, im, maxIter) as PointInfo,
  defaultBounds: defaultBoundsJs,
  generateMusic: generateMusicJs,
};

let loadPromise: Promise<FractalsApi> | null = null;

async function tryLoadWasm(assetPath: string): Promise<FractalsApi | null> {
  try {
    const jsUrl = new URL(assetPath.replace(/\/?$/, "/") + "fractals_wasm.js", window.location.href)
      .toString();
    const mod = await import(/* @vite-ignore */ /* webpackIgnore: true */ jsUrl);
    // wasm-pack's web target expects an init() call with the wasm URL.
    const wasmUrl = new URL(
      assetPath.replace(/\/?$/, "/") + "fractals_wasm_bg.wasm",
      window.location.href
    ).toString();
    const init = mod.default ?? mod.init;
    if (typeof init !== "function") return null;
    await init(wasmUrl);
    return {
      kind: "wasm",
      render: (kind, cRe, cIm, xMin, xMax, yMin, yMax, width, height, maxIter) =>
        mod.render_fractal(
          kind,
          cRe,
          cIm,
          xMin,
          xMax,
          yMin,
          yMax,
          width,
          height,
          maxIter
        ) as Float32Array,
      renderProgressive: (kind, cRe, cIm, xMin, xMax, yMin, yMax, width, height, maxIter, stride) =>
        mod.render_fractal_progressive(
          kind,
          cRe,
          cIm,
          xMin,
          xMax,
          yMin,
          yMax,
          width,
          height,
          maxIter,
          stride
        ) as Float32Array,
      computePoint: (kind, cRe, cIm, re, im, maxIter) =>
        mod.compute_point(kind, cRe, cIm, re, im, maxIter) as PointInfo,
      defaultBounds: (kind) => mod.default_bounds(kind),
      generateMusic: (cfg) => mod.generate_music(cfg) as GeneratedNote[],
    };
  } catch (err) {
    // Intentionally swallow — we'll fall back to JS.
    // eslint-disable-next-line no-console
    console.warn("[fractals/explorer] WASM load failed, falling back to JS.", err);
    return null;
  }
}

/**
 * Load the fractal compute API. Returns the WASM-backed version when
 * available, otherwise the JS fallback. Safe to call multiple times — the
 * result is memoized per module load.
 */
export function loadFractalsApi(assetPath = "./wasm/"): Promise<FractalsApi> {
  if (typeof window === "undefined") {
    // SSR: hand back the JS fallback immediately.
    return Promise.resolve(jsApi);
  }
  if (!loadPromise) {
    loadPromise = tryLoadWasm(assetPath).then((api) => api ?? jsApi);
  }
  return loadPromise;
}
