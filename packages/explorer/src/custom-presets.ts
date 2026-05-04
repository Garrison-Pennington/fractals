import type { CustomFractalDef } from "./types.js";

export const CUSTOM_PRESETS: CustomFractalDef[] = [
  {
    name: "Mandelbrot (z^2 + c)",
    initZr: "0",
    initZi: "0",
    stepZr: "zr*zr - zi*zi + cr",
    stepZi: "2*zr*zi + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -2.5, x_max: 1.0, y_min: -1.25, y_max: 1.25 },
  },
  {
    name: "Tricorn (Mandelbar)",
    initZr: "0",
    initZi: "0",
    stepZr: "zr*zr - zi*zi + cr",
    stepZi: "-2*zr*zi + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -2.5, x_max: 1.5, y_min: -1.5, y_max: 1.5 },
  },
  {
    name: "Cubic Mandelbrot (z^3 + c)",
    initZr: "0",
    initZi: "0",
    stepZr: "zr*zr*zr - 3*zr*zi*zi + cr",
    stepZi: "3*zr*zr*zi - zi*zi*zi + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -1.6, x_max: 1.6, y_min: -1.2, y_max: 1.2 },
  },
  {
    name: "Quartic Mandelbrot (z^4 + c)",
    initZr: "0",
    initZi: "0",
    stepZr: "zr*zr*zr*zr - 6*zr*zr*zi*zi + zi*zi*zi*zi + cr",
    stepZi: "4*zr*zr*zr*zi - 4*zr*zi*zi*zi + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -1.5, x_max: 1.5, y_min: -1.2, y_max: 1.2 },
  },
  {
    name: "Burning Ship",
    initZr: "0",
    initZi: "0",
    stepZr: "zr*zr - zi*zi + cr",
    stepZi: "2*Math.abs(zr*zi) + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -2.5, x_max: 1.5, y_min: -2.0, y_max: 1.0 },
  },
  {
    name: "Celtic",
    initZr: "0",
    initZi: "0",
    stepZr: "Math.abs(zr*zr - zi*zi) + cr",
    stepZi: "2*zr*zi + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -2.5, x_max: 1.5, y_min: -1.5, y_max: 1.5 },
  },
  {
    name: "Buffalo",
    initZr: "0",
    initZi: "0",
    stepZr: "Math.abs(zr*zr - zi*zi) + cr",
    stepZi: "Math.abs(2*zr*zi) + ci",
    escapeExpr: "zr*zr + zi*zi > 4",
    defaultBounds: { x_min: -2.5, x_max: 1.5, y_min: -1.5, y_max: 1.5 },
  },
];

export const DEFAULT_CUSTOM_FRACTAL: CustomFractalDef = CUSTOM_PRESETS[0];
