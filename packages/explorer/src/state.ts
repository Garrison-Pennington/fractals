// Shareable-URL state encoding.
//
// The full ExplorerState is serialised to JSON, deflate-compressed, and
// base64-URL encoded. The first byte of the compressed payload is a version
// marker so old share links remain readable after the format evolves.

import { deflateSync, inflateSync } from "fflate";
import type { CustomFractalDef, ExplorerState, FractalKind, JuliaKeyframe, PaletteName, ScaleKind, SequenceType, Waveform } from "./types.js";
import { DEFAULT_CUSTOM_FRACTAL } from "./custom-presets.js";

const VERSION: number = 1;

export const DEFAULT_STATE: ExplorerState = {
  fractal: "mandelbrot",
  bounds: { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 },
  maxIterations: 400,
  palette: "classic",
  julia: { re: -0.7, im: 0.27015 },
  music: {
    scale: "minor_pentatonic",
    rootMidi: 60,
    sequenceType: "path",
    length: 16,
    waveform: "triangle",
  },
  keyframes: [],
  customFractal: DEFAULT_CUSTOM_FRACTAL,
};

export function defaultsForFractal(kind: FractalKind, customDef?: CustomFractalDef) {
  switch (kind) {
    case "julia":
      return { x_min: -1.5, x_max: 1.5, y_min: -1.5, y_max: 1.5 };
    case "burning_ship":
      return { x_min: -2.2, x_max: 1.5, y_min: -2.0, y_max: 1.2 };
    case "newton":
      return { x_min: -2.0, x_max: 2.0, y_min: -2.0, y_max: 2.0 };
    case "custom":
      if (customDef) return { ...customDef.defaultBounds };
      return { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 };
    default:
      return { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 };
  }
}

function toBase64Url(bytes: Uint8Array): string {
  let bin = "";
  for (let i = 0; i < bytes.length; i++) bin += String.fromCharCode(bytes[i]);
  return btoa(bin).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function fromBase64Url(s: string): Uint8Array {
  const pad = "=".repeat((4 - (s.length % 4)) % 4);
  const b64 = (s + pad).replace(/-/g, "+").replace(/_/g, "/");
  const bin = atob(b64);
  const out = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) out[i] = bin.charCodeAt(i);
  return out;
}

export function encodeState(state: ExplorerState): string {
  const json = JSON.stringify(state);
  const data = new TextEncoder().encode(json);
  const compressed = deflateSync(data);
  const withVersion = new Uint8Array(compressed.length + 1);
  withVersion[0] = VERSION;
  withVersion.set(compressed, 1);
  return toBase64Url(withVersion);
}

export function decodeState(encoded: string): ExplorerState | null {
  try {
    const bytes = fromBase64Url(encoded);
    if (bytes.length < 2) return null;
    const version = bytes[0];
    if (version !== 1) {
      console.warn("[fractals/explorer] unknown state version", version);
      return null;
    }
    const payload = bytes.slice(1);
    const json = new TextDecoder().decode(inflateSync(payload));
    const parsed = JSON.parse(json);
    return coerceState(parsed);
  } catch (err) {
    console.warn("[fractals/explorer] state decode failed", err);
    return null;
  }
}

function coerceCustomDef(raw: unknown): CustomFractalDef | null {
  if (!raw || typeof raw !== "object") return null;
  const r = raw as Record<string, unknown>;
  if (
    typeof r.initZr !== "string" ||
    typeof r.initZi !== "string" ||
    typeof r.stepZr !== "string" ||
    typeof r.stepZi !== "string" ||
    typeof r.escapeExpr !== "string"
  )
    return null;
  const bounds = r.defaultBounds as CustomFractalDef["defaultBounds"] | undefined;
  return {
    name: typeof r.name === "string" ? r.name : "Custom",
    initZr: r.initZr,
    initZi: r.initZi,
    stepZr: r.stepZr,
    stepZi: r.stepZi,
    escapeExpr: r.escapeExpr,
    defaultBounds: bounds && typeof bounds.x_min === "number"
      ? bounds
      : { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 },
  };
}

function coerceState(raw: unknown): ExplorerState | null {
  if (!raw || typeof raw !== "object") return null;
  const r = raw as Record<string, unknown>;
  const fractal = (r.fractal as FractalKind) || "mandelbrot";
  const bounds = r.bounds as ExplorerState["bounds"];
  if (
    !bounds ||
    typeof bounds.x_min !== "number" ||
    typeof bounds.x_max !== "number" ||
    typeof bounds.y_min !== "number" ||
    typeof bounds.y_max !== "number"
  ) {
    return null;
  }
  const julia = (r.julia as ExplorerState["julia"]) || { re: -0.7, im: 0.27015 };
  const music = r.music as ExplorerState["music"] | undefined;
  const keyframes = Array.isArray(r.keyframes) ? (r.keyframes as JuliaKeyframe[]) : [];
  return {
    fractal: ["mandelbrot", "julia", "burning_ship", "newton", "custom"].includes(fractal)
      ? fractal
      : "mandelbrot",
    bounds,
    maxIterations: typeof r.maxIterations === "number" ? r.maxIterations : 400,
    palette: (["classic", "fire", "rainbow"] as PaletteName[]).includes(
      r.palette as PaletteName
    )
      ? (r.palette as PaletteName)
      : "classic",
    julia: { re: Number(julia.re) || 0, im: Number(julia.im) || 0 },
    music: {
      scale: (music?.scale as ScaleKind) || "minor_pentatonic",
      rootMidi: typeof music?.rootMidi === "number" ? music.rootMidi : 60,
      sequenceType: (music?.sequenceType as SequenceType) || "path",
      length: typeof music?.length === "number" ? music.length : 16,
      waveform: (music?.waveform as Waveform) || "triangle",
    },
    keyframes: keyframes.filter(
      (k): k is JuliaKeyframe =>
        !!k && typeof k.c?.re === "number" && typeof k.c?.im === "number" && typeof k.duration === "number"
    ),
    customFractal: coerceCustomDef(r.customFractal) ?? DEFAULT_CUSTOM_FRACTAL,
  };
}

export function encodeStateToHash(state: ExplorerState): string {
  return `#state=${encodeState(state)}`;
}

export function readStateFromHash(): ExplorerState | null {
  if (typeof window === "undefined") return null;
  const hash = window.location.hash;
  const match = hash.match(/state=([^&]+)/);
  if (!match) return null;
  return decodeState(match[1]);
}
