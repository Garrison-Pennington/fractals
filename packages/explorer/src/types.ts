// Shared type definitions for the fractal explorer.

export type FractalKind = "mandelbrot" | "julia" | "burning_ship" | "newton" | "custom";

export interface Bounds {
  x_min: number;
  x_max: number;
  y_min: number;
  y_max: number;
}

export interface JuliaParams {
  re: number;
  im: number;
}

export interface JuliaKeyframe {
  c: JuliaParams;
  // Seconds to reach this keyframe from the previous one.
  duration: number;
}

export type PaletteName = "classic" | "fire" | "rainbow";
export type ScaleKind = "major" | "minor" | "minor_pentatonic" | "major_pentatonic" | "blues";
export type SequenceType = "path" | "region";
export type Waveform = "sine" | "triangle" | "sawtooth";

export interface MusicConfig {
  scale: ScaleKind;
  rootMidi: number;
  sequenceType: SequenceType;
  length: number;
  waveform: Waveform;
}

export interface GeneratedNote {
  pitch: number;
  duration: number;
  velocity: number;
  name: string;
}

/** User-defined fractal: expressions are JS snippets evaluated per pixel. */
export interface CustomFractalDef {
  name: string;
  /** Initial zr expression.  Available vars: re, im (pixel coords), cr, ci (c param). */
  initZr: string;
  /** Initial zi expression. */
  initZi: string;
  /** Next zr expression.  Available vars: zr, zi, cr, ci, re, im. */
  stepZr: string;
  /** Next zi expression. */
  stepZi: string;
  /** Boolean escape expression, e.g. "zr*zr + zi*zi > 4". */
  escapeExpr: string;
  defaultBounds: Bounds;
}

export interface ExplorerState {
  fractal: FractalKind;
  bounds: Bounds;
  maxIterations: number;
  palette: PaletteName;
  julia: JuliaParams;
  music: MusicConfig;
  keyframes: JuliaKeyframe[];
  customFractal: CustomFractalDef;
}

export interface PointInfo {
  re: number;
  im: number;
  iterations: number;
  escaped: boolean;
  smooth: number;
}
