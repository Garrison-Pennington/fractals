"use client";
// FractalExplorer — the full interactive component.
//
// Feature layout:
//   [viewport canvas + coordinate status bar]
//   [controls panel: fractal type, Julia c sliders, iterations, palette, zoom back, share]
//   [Julia animation panel: keyframe editor + play/record]
//   [Music sequencer panel: scale/root/sequence config + piano roll + play/stop]
//
// Rendering strategy: we always have the JS fallback ready synchronously (so
// the first paint happens before WASM init resolves). Once WASM loads, we
// swap in the faster backend. Every render pass runs at 1/8 -> 1/4 -> 1/2 ->
// 1/1 stride so users see the image sharpen from the first frame.

import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  CSSProperties,
} from "react";
import JSZip from "jszip";
import { colorize, PALETTE_NAMES } from "./color.js";
import { loadFractalsApi, jsApi, setActiveCustomDef, FractalsApi } from "./wasm-loader.js";
import { compileCustomFractal } from "./fallback.js";
import { CUSTOM_PRESETS } from "./custom-presets.js";
import {
  DEFAULT_STATE,
  defaultsForFractal,
  encodeStateToHash,
  readStateFromHash,
} from "./state.js";
import { playSequence, PlayHandle } from "./audio.js";
import type {
  Bounds,
  CustomFractalDef,
  ExplorerState,
  FractalKind,
  GeneratedNote,
  JuliaKeyframe,
  PaletteName,
  PointInfo,
  ScaleKind,
  SequenceType,
  Waveform,
} from "./types.js";

export interface FractalExplorerProps {
  /** Width in CSS pixels. Defaults to 100% of parent. */
  width?: number | string;
  /** Height in CSS pixels. Defaults to 100% of parent. */
  height?: number | string;
  /** Path from which to load wasm assets. Defaults to "./wasm/". */
  assetPath?: string;
  className?: string;
  style?: CSSProperties;
}

interface ProgressivePass {
  stride: number;
}

const PASSES: ProgressivePass[] = [{ stride: 8 }, { stride: 4 }, { stride: 2 }, { stride: 1 }];
const SCALE_OPTIONS: ScaleKind[] = ["major", "minor", "major_pentatonic", "minor_pentatonic", "blues"];
const WAVEFORM_OPTIONS: Waveform[] = ["sine", "triangle", "sawtooth"];
const FRACTAL_OPTIONS: FractalKind[] = ["mandelbrot", "julia", "burning_ship", "newton", "custom"];

function cubicEaseInOut(t: number): number {
  return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
}

function formatComplex(re: number, im: number, prec = 6): string {
  const s = im >= 0 ? "+" : "-";
  return `${re.toFixed(prec)} ${s} ${Math.abs(im).toFixed(prec)}i`;
}

export const FractalExplorer: React.FC<FractalExplorerProps> = ({
  width = "100%",
  height = "100%",
  assetPath = "./wasm/",
  className,
  style,
}) => {
  // ------ Core state ------
  const [state, setState] = useState<ExplorerState>(DEFAULT_STATE);
  const [api, setApi] = useState<FractalsApi>(jsApi);
  const [pointInfo, setPointInfo] = useState<PointInfo | null>(null);
  const [notes, setNotes] = useState<GeneratedNote[]>([]);
  const [playingIdx, setPlayingIdx] = useState<number>(-1);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isAnimating, setIsAnimating] = useState(false);
  const [recordingFrames, setRecordingFrames] = useState(false);
  const [shareStatus, setShareStatus] = useState<string>("");
  const [resolution, setResolution] = useState<{ w: number; h: number }>({ w: 640, h: 480 });
  const historyRef = useRef<ExplorerState[]>([]);

  const containerRef = useRef<HTMLDivElement | null>(null);
  const canvasWrapperRef = useRef<HTMLDivElement | null>(null);
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const selectionRef = useRef<{ start: [number, number]; end: [number, number] } | null>(null);
  const [selectionBox, setSelectionBox] = useState<{
    x: number;
    y: number;
    w: number;
    h: number;
  } | null>(null);
  const animationRef = useRef<number | null>(null);
  const audioHandleRef = useRef<PlayHandle | null>(null);
  const currentRenderIdRef = useRef<number>(0);
  const lastInteractionRef = useRef<number>(Date.now());
  const frameCacheRef = useRef<Blob[]>([]);
  const lastFrameTimeRef = useRef<number>(0);
  const adaptiveStrideRef = useRef<number>(1);

  // ------ Initial state from URL + WASM load ------
  useEffect(() => {
    // Read shared state, if any, before loading WASM.
    const fromHash = readStateFromHash();
    if (fromHash) setState(fromHash);
    let disposed = false;
    loadFractalsApi(assetPath).then((a) => {
      if (!disposed) setApi(a);
    });
    return () => {
      disposed = true;
    };
  }, [assetPath]);

  // ------ Resize observer: canvas matches its wrapper (not the outer container,
  //        which would create a feedback loop via panel height changes). ------
  useEffect(() => {
    const el = canvasWrapperRef.current;
    if (!el) return;
    const ro = new ResizeObserver(() => {
      const rect = el.getBoundingClientRect();
      const w = Math.max(320, Math.round(rect.width));
      const h = Math.max(240, Math.round(rect.height));
      setResolution({ w, h });
    });
    ro.observe(el);
    return () => ro.disconnect();
  }, []);

  // ------ Update location hash when state settles ------
  useEffect(() => {
    if (typeof window === "undefined") return;
    const hash = encodeStateToHash(state);
    // Replace instead of push so we don't spam the back button.
    if (window.location.hash !== hash) {
      window.history.replaceState(null, "", `${window.location.pathname}${window.location.search}${hash}`);
    }
  }, [state]);

  // ------ Core render loop: triggered by state / resolution / api change ------
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    const width = resolution.w;
    const height = resolution.h;
    canvas.width = width;
    canvas.height = height;

    const id = ++currentRenderIdRef.current;
    const isNewton = state.fractal === "newton";
    const cRe = state.julia.re;
    const cIm = state.julia.im;

    // Custom fractals run through the JS fallback (WASM doesn't know them).
    const isCustom = state.fractal === "custom";
    if (isCustom) setActiveCustomDef(state.customFractal);
    const renderApi = isCustom ? jsApi : api;

    const run = async () => {
      // Stride sequence: 8 -> 4 -> 2 -> 1. Early exit if state changes.
      for (const { stride } of PASSES) {
        if (id !== currentRenderIdRef.current) return;
        // Yield to paint between passes so the previous frame is visible.
        await new Promise<void>((r) => requestAnimationFrame(() => r()));
        if (id !== currentRenderIdRef.current) return;

        let buf: Float32Array;
        try {
          buf =
            stride === 1
              ? renderApi.render(
                  state.fractal,
                  cRe,
                  cIm,
                  state.bounds.x_min,
                  state.bounds.x_max,
                  state.bounds.y_min,
                  state.bounds.y_max,
                  width,
                  height,
                  state.maxIterations
                )
              : renderApi.renderProgressive(
                  state.fractal,
                  cRe,
                  cIm,
                  state.bounds.x_min,
                  state.bounds.x_max,
                  state.bounds.y_min,
                  state.bounds.y_max,
                  width,
                  height,
                  state.maxIterations,
                  stride
                );
        } catch {
          return; // bad custom expression — stop rendering
        }

        if (id !== currentRenderIdRef.current) return;
        const rgba = colorize(buf, width, height, state.palette, isNewton);
        // Safari + TypeScript are finicky about Uint8ClampedArray buffer types;
        // write through an ImageData obtained from the context to sidestep it.
        const imageData = ctx.createImageData(width, height);
        imageData.data.set(rgba);
        ctx.putImageData(imageData, 0, 0);
      }
    };
    run();
  }, [state, api, resolution]);

  // ------ Pointer handling: drag-to-zoom, hover, wheel ------
  const toComplex = useCallback(
    (clientX: number, clientY: number): [number, number] => {
      const canvas = canvasRef.current!;
      const rect = canvas.getBoundingClientRect();
      const u = (clientX - rect.left) / rect.width;
      const v = (clientY - rect.top) / rect.height;
      const x = state.bounds.x_min + u * (state.bounds.x_max - state.bounds.x_min);
      const y = state.bounds.y_min + v * (state.bounds.y_max - state.bounds.y_min);
      return [x, y];
    },
    [state.bounds]
  );

  const handlePointerDown = (e: React.PointerEvent<HTMLCanvasElement>) => {
    const canvas = canvasRef.current!;
    canvas.setPointerCapture(e.pointerId);
    const rect = canvas.getBoundingClientRect();
    selectionRef.current = {
      start: [e.clientX - rect.left, e.clientY - rect.top],
      end: [e.clientX - rect.left, e.clientY - rect.top],
    };
    setSelectionBox({ x: e.clientX - rect.left, y: e.clientY - rect.top, w: 0, h: 0 });
  };
  const handlePointerMove = (e: React.PointerEvent<HTMLCanvasElement>) => {
    lastInteractionRef.current = Date.now();
    const [re, im] = toComplex(e.clientX, e.clientY);
    // Avoid running compute for every single move if WASM hasn't loaded yet.
    const pointApi = state.fractal === "custom" ? jsApi : api;
    const info = pointApi.computePoint(state.fractal, state.julia.re, state.julia.im, re, im, 64);
    setPointInfo({ ...info, re, im });
    if (selectionRef.current) {
      const canvas = canvasRef.current!;
      const rect = canvas.getBoundingClientRect();
      selectionRef.current.end = [e.clientX - rect.left, e.clientY - rect.top];
      const [sx, sy] = selectionRef.current.start;
      const [ex, ey] = selectionRef.current.end;
      setSelectionBox({
        x: Math.min(sx, ex),
        y: Math.min(sy, ey),
        w: Math.abs(ex - sx),
        h: Math.abs(ey - sy),
      });
    }
  };
  const handlePointerUp = () => {
    const sel = selectionRef.current;
    selectionRef.current = null;
    setSelectionBox(null);
    if (!sel) return;
    const [sx, sy] = sel.start;
    const [ex, ey] = sel.end;
    if (Math.abs(ex - sx) < 5 || Math.abs(ey - sy) < 5) return;
    const canvas = canvasRef.current!;
    const rect = canvas.getBoundingClientRect();
    const x1 = Math.min(sx, ex) / rect.width;
    const x2 = Math.max(sx, ex) / rect.width;
    const y1 = Math.min(sy, ey) / rect.height;
    const y2 = Math.max(sy, ey) / rect.height;
    const nb: Bounds = {
      x_min: state.bounds.x_min + x1 * (state.bounds.x_max - state.bounds.x_min),
      x_max: state.bounds.x_min + x2 * (state.bounds.x_max - state.bounds.x_min),
      y_min: state.bounds.y_min + y1 * (state.bounds.y_max - state.bounds.y_min),
      y_max: state.bounds.y_min + y2 * (state.bounds.y_max - state.bounds.y_min),
    };
    historyRef.current.push(state);
    setState({ ...state, bounds: nb });
  };

  const handleWheel = useCallback((e: WheelEvent) => {
    e.preventDefault();
    const [re, im] = toComplex(e.clientX, e.clientY);
    const factor = e.deltaY < 0 ? 1.25 : 1 / 1.25;
    setState((s) => {
      historyRef.current.push(s);
      const w = (s.bounds.x_max - s.bounds.x_min) / factor;
      const h = (s.bounds.y_max - s.bounds.y_min) / factor;
      return {
        ...s,
        bounds: {
          x_min: re - w / 2,
          x_max: re + w / 2,
          y_min: im - h / 2,
          y_max: im + h / 2,
        },
      };
    });
  }, [toComplex]);

  // Attach wheel listener as non-passive so preventDefault works.
  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    canvas.addEventListener("wheel", handleWheel, { passive: false });
    return () => canvas.removeEventListener("wheel", handleWheel);
  }, [handleWheel]);

  const zoomBack = () => {
    const prev = historyRef.current.pop();
    if (prev) setState(prev);
  };
  const resetView = () => {
    historyRef.current.push(state);
    setState({ ...state, bounds: defaultsForFractal(state.fractal, state.customFractal) });
  };

  // ------ Fractal type switch ------
  const onFractalChange = (kind: FractalKind) => {
    historyRef.current.push(state);
    setState({ ...state, fractal: kind, bounds: defaultsForFractal(kind, state.customFractal) });
  };

  // ------ Julia animation ------
  const addKeyframe = () => {
    const k: JuliaKeyframe = { c: { ...state.julia }, duration: 1.0 };
    setState({ ...state, keyframes: [...state.keyframes, k] });
  };
  const removeKeyframe = (i: number) => {
    setState({ ...state, keyframes: state.keyframes.filter((_, idx) => idx !== i) });
  };
  const updateKeyframeDuration = (i: number, duration: number) => {
    const kf = state.keyframes.slice();
    kf[i] = { ...kf[i], duration };
    setState({ ...state, keyframes: kf });
  };

  const playAnimation = useCallback(() => {
    if (state.keyframes.length < 2) return;
    setIsAnimating(true);
    const start = performance.now();
    const total = state.keyframes.reduce((a, k) => a + k.duration, 0) * 1000;
    const step = (now: number) => {
      const elapsed = (now - start) % total;
      let acc = 0;
      let fromIdx = 0;
      for (let i = 0; i < state.keyframes.length; i++) {
        const dur = state.keyframes[i].duration * 1000;
        if (elapsed < acc + dur) {
          fromIdx = i;
          break;
        }
        acc += dur;
      }
      const toIdx = (fromIdx + 1) % state.keyframes.length;
      const localDur = state.keyframes[fromIdx].duration * 1000;
      const t = localDur > 0 ? cubicEaseInOut(Math.min(1, (elapsed - acc) / localDur)) : 0;
      const from = state.keyframes[fromIdx].c;
      const to = state.keyframes[toIdx].c;
      const nextJulia = {
        re: from.re + (to.re - from.re) * t,
        im: from.im + (to.im - from.im) * t,
      };
      // Adaptive resolution: if last frame took too long, drop stride for the
      // next. We do this by adjusting maxIterations on the fly for animation
      // playback — the full state isn't mutated, just a transient julia.re/im.
      setState((s) => ({ ...s, julia: nextJulia }));
      animationRef.current = requestAnimationFrame(step);
      // Recording: capture current frame.
      if (recordingFrames && canvasRef.current) {
        const framesPerSec = 20;
        if (now - lastFrameTimeRef.current > 1000 / framesPerSec) {
          lastFrameTimeRef.current = now;
          canvasRef.current.toBlob((blob) => {
            if (blob) frameCacheRef.current.push(blob);
          }, "image/png");
        }
      }
    };
    animationRef.current = requestAnimationFrame(step);
  }, [state.keyframes, recordingFrames]);

  const stopAnimation = useCallback(() => {
    if (animationRef.current != null) cancelAnimationFrame(animationRef.current);
    animationRef.current = null;
    setIsAnimating(false);
    if (recordingFrames) {
      finalizeRecording();
      setRecordingFrames(false);
    }
  }, [recordingFrames]);

  async function finalizeRecording() {
    if (frameCacheRef.current.length === 0) return;
    const zip = new JSZip();
    for (let i = 0; i < frameCacheRef.current.length; i++) {
      const blob = frameCacheRef.current[i];
      zip.file(`frame_${String(i).padStart(5, "0")}.png`, blob);
    }
    const out = await zip.generateAsync({ type: "blob" });
    const a = document.createElement("a");
    a.href = URL.createObjectURL(out);
    a.download = "fractal_animation.zip";
    a.click();
    frameCacheRef.current = [];
  }

  // ------ Music ------
  const generateMusic = () => {
    const out = api.generateMusic({
      kind: state.fractal,
      c_re: state.julia.re,
      c_im: state.julia.im,
      scale: state.music.scale,
      root_midi: state.music.rootMidi,
      sequence_type: state.music.sequenceType,
      length: state.music.length,
      max_iter: state.maxIterations,
      x_min: state.bounds.x_min,
      x_max: state.bounds.x_max,
      y_min: state.bounds.y_min,
      y_max: state.bounds.y_max,
    });
    setNotes(out);
  };

  const playMusic = () => {
    if (notes.length === 0) generateMusic();
    const n = notes.length === 0 ? api.generateMusic({
      kind: state.fractal,
      c_re: state.julia.re,
      c_im: state.julia.im,
      scale: state.music.scale,
      root_midi: state.music.rootMidi,
      sequence_type: state.music.sequenceType,
      length: state.music.length,
      max_iter: state.maxIterations,
      x_min: state.bounds.x_min,
      x_max: state.bounds.x_max,
      y_min: state.bounds.y_min,
      y_max: state.bounds.y_max,
    }) : notes;
    if (notes.length === 0) setNotes(n);
    audioHandleRef.current?.stop();
    setIsPlaying(true);
    setPlayingIdx(-1);
    audioHandleRef.current = playSequence(n, {
      waveform: state.music.waveform,
      bpm: 120,
      onNote: (i) => setPlayingIdx(i),
      onEnd: () => {
        setIsPlaying(false);
        setPlayingIdx(-1);
      },
    });
  };

  const stopMusic = () => {
    audioHandleRef.current?.stop();
    audioHandleRef.current = null;
    setIsPlaying(false);
    setPlayingIdx(-1);
  };

  // ------ Share URL ------
  const share = async () => {
    const url = `${window.location.origin}${window.location.pathname}${window.location.search}${encodeStateToHash(state)}`;
    try {
      await navigator.clipboard.writeText(url);
      setShareStatus("Link copied!");
    } catch {
      setShareStatus("Copy failed — URL is in the address bar");
    }
    window.setTimeout(() => setShareStatus(""), 2500);
  };

  // ------ Cleanup ------
  useEffect(() => {
    return () => {
      audioHandleRef.current?.stop();
      if (animationRef.current != null) cancelAnimationFrame(animationRef.current);
    };
  }, []);

  // ------ Render ------
  const containerStyle: CSSProperties = {
    width,
    height,
    display: "flex",
    flexDirection: "column",
    overflow: "hidden",
    minHeight: 0,
    background: "#0b0b12",
    color: "#e8e8ee",
    fontFamily: "system-ui, -apple-system, Segoe UI, sans-serif",
    fontSize: 13,
    ...style,
  };

  return (
    <div ref={containerRef} className={className} style={containerStyle}>
      <div ref={canvasWrapperRef} style={{ position: "relative", flex: 1, minHeight: 0 }}>
        <canvas
          ref={canvasRef}
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerUp}
          onPointerCancel={handlePointerUp}
          style={{
            width: "100%",
            height: "100%",
            display: "block",
            cursor: selectionBox ? "crosshair" : "grab",
            imageRendering: "auto",
            touchAction: "none",
          }}
        />
        {selectionBox && (
          <div
            style={{
              position: "absolute",
              left: selectionBox.x,
              top: selectionBox.y,
              width: selectionBox.w,
              height: selectionBox.h,
              border: "1.5px dashed #fff",
              background: "rgba(255,255,255,0.08)",
              pointerEvents: "none",
            }}
          />
        )}
        <div
          style={{
            position: "absolute",
            left: 8,
            top: 8,
            background: "rgba(10,10,18,0.75)",
            padding: "6px 10px",
            borderRadius: 6,
            fontFamily: "ui-monospace, SFMono-Regular, monospace",
            fontSize: 11,
            pointerEvents: "none",
          }}
        >
          <div>
            bounds: [{state.bounds.x_min.toExponential(3)}, {state.bounds.x_max.toExponential(3)}] ×
            [{state.bounds.y_min.toExponential(3)}, {state.bounds.y_max.toExponential(3)}]
          </div>
          <div>
            zoom: {(4 / (state.bounds.x_max - state.bounds.x_min)).toExponential(2)} · backend:{" "}
            {api.kind.toUpperCase()}
          </div>
        </div>
        <div
          style={{
            position: "absolute",
            right: 8,
            top: 8,
            background: "rgba(10,10,18,0.75)",
            padding: "6px 10px",
            borderRadius: 6,
            fontFamily: "ui-monospace, SFMono-Regular, monospace",
            fontSize: 11,
            pointerEvents: "none",
          }}
        >
          {pointInfo
            ? `${formatComplex(pointInfo.re, pointInfo.im, 5)}  iter=${pointInfo.iterations}${pointInfo.escaped ? "" : " ∞"}`
            : "hover to inspect"}
        </div>
      </div>

      <ControlsPanel
        state={state}
        setState={setState}
        onFractalChange={onFractalChange}
        zoomBack={zoomBack}
        resetView={resetView}
        share={share}
        shareStatus={shareStatus}
        canUndo={historyRef.current.length > 0}
      />

      {state.fractal === "custom" && (
        <CustomFractalPanel state={state} setState={setState} />
      )}

      <AnimationPanel
        state={state}
        addKeyframe={addKeyframe}
        removeKeyframe={removeKeyframe}
        updateDuration={updateKeyframeDuration}
        isAnimating={isAnimating}
        onPlay={playAnimation}
        onStop={stopAnimation}
        recording={recordingFrames}
        toggleRecording={() => setRecordingFrames((r) => !r)}
      />

      <MusicPanel
        state={state}
        setState={setState}
        notes={notes}
        playingIdx={playingIdx}
        isPlaying={isPlaying}
        onGenerate={generateMusic}
        onPlay={playMusic}
        onStop={stopMusic}
      />
    </div>
  );
};

// ----------------------------------------------------------------------------
// Controls Panel
// ----------------------------------------------------------------------------
const panelStyle: CSSProperties = {
  display: "flex",
  flexWrap: "wrap",
  gap: 12,
  padding: "8px 12px",
  background: "#11111a",
  borderTop: "1px solid #22222c",
  alignItems: "center",
};
const labelStyle: CSSProperties = {
  display: "flex",
  flexDirection: "column",
  gap: 2,
  fontSize: 11,
  color: "#a0a0b0",
};
const inputStyle: CSSProperties = {
  background: "#1a1a25",
  color: "#fff",
  border: "1px solid #2c2c3a",
  borderRadius: 4,
  padding: "3px 6px",
  fontSize: 12,
};
const buttonStyle: CSSProperties = {
  background: "#2c3144",
  color: "#fff",
  border: "1px solid #383e54",
  borderRadius: 4,
  padding: "4px 10px",
  fontSize: 12,
  cursor: "pointer",
};

const ControlsPanel: React.FC<{
  state: ExplorerState;
  setState: React.Dispatch<React.SetStateAction<ExplorerState>>;
  onFractalChange: (kind: FractalKind) => void;
  zoomBack: () => void;
  resetView: () => void;
  share: () => void;
  shareStatus: string;
  canUndo: boolean;
}> = ({ state, setState, onFractalChange, zoomBack, resetView, share, shareStatus, canUndo }) => {
  return (
    <div style={panelStyle}>
      <label style={labelStyle}>
        Fractal
        <select
          style={inputStyle}
          value={state.fractal}
          onChange={(e) => onFractalChange(e.target.value as FractalKind)}
        >
          {FRACTAL_OPTIONS.map((k) => (
            <option key={k} value={k}>
              {k.replace("_", " ")}
            </option>
          ))}
        </select>
      </label>

      {state.fractal === "julia" && (
        <>
          <label style={labelStyle}>
            c.re: {state.julia.re.toFixed(4)}
            <input
              type="range"
              min={-2}
              max={2}
              step={0.001}
              value={state.julia.re}
              onChange={(e) =>
                setState({ ...state, julia: { ...state.julia, re: parseFloat(e.target.value) } })
              }
            />
          </label>
          <label style={labelStyle}>
            c.im: {state.julia.im.toFixed(4)}
            <input
              type="range"
              min={-2}
              max={2}
              step={0.001}
              value={state.julia.im}
              onChange={(e) =>
                setState({ ...state, julia: { ...state.julia, im: parseFloat(e.target.value) } })
              }
            />
          </label>
        </>
      )}

      <label style={labelStyle}>
        max iter: {state.maxIterations}
        <input
          type="range"
          min={50}
          max={5000}
          step={50}
          value={state.maxIterations}
          onChange={(e) => setState({ ...state, maxIterations: parseInt(e.target.value) })}
        />
      </label>

      <label style={labelStyle}>
        palette
        <select
          style={inputStyle}
          value={state.palette}
          onChange={(e) => setState({ ...state, palette: e.target.value as PaletteName })}
        >
          {PALETTE_NAMES.map((p) => (
            <option key={p} value={p}>
              {p}
            </option>
          ))}
        </select>
      </label>

      <button style={buttonStyle} onClick={zoomBack} disabled={!canUndo}>
        ← Back
      </button>
      <button style={buttonStyle} onClick={resetView}>
        Reset view
      </button>
      <button style={buttonStyle} onClick={share}>
        Share
      </button>
      {shareStatus && <span style={{ color: "#8afc9a", fontSize: 11 }}>{shareStatus}</span>}
    </div>
  );
};

// ----------------------------------------------------------------------------
// Custom Fractal Panel
// ----------------------------------------------------------------------------
const exprInputStyle: CSSProperties = {
  ...inputStyle,
  fontFamily: "ui-monospace, SFMono-Regular, monospace",
  fontSize: 11,
  width: "100%",
  boxSizing: "border-box",
};

const CustomFractalPanel: React.FC<{
  state: ExplorerState;
  setState: React.Dispatch<React.SetStateAction<ExplorerState>>;
}> = ({ state, setState }) => {
  const def = state.customFractal;
  const [error, setError] = useState<string | null>(null);

  const updateField = (field: keyof CustomFractalDef, value: string) => {
    const next = { ...def, [field]: value, name: "Custom" };
    // Validate by attempting compilation.
    try {
      compileCustomFractal(next);
      setError(null);
    } catch (e: unknown) {
      setError((e as Error).message);
    }
    setState((s) => ({ ...s, customFractal: next }));
  };

  const loadPreset = (preset: CustomFractalDef) => {
    setError(null);
    setState((s) => ({
      ...s,
      customFractal: preset,
      bounds: { ...preset.defaultBounds },
    }));
  };

  return (
    <div style={{ ...panelStyle, flexDirection: "column", alignItems: "stretch", gap: 6 }}>
      <div style={{ display: "flex", gap: 8, alignItems: "center", flexWrap: "wrap" }}>
        <span style={{ fontSize: 11, color: "#aab", minWidth: 50 }}>Custom fractal</span>
        <select
          style={{ ...inputStyle, flex: 1, minWidth: 140 }}
          value={CUSTOM_PRESETS.find(
            (p) => p.stepZr === def.stepZr && p.stepZi === def.stepZi && p.initZr === def.initZr
          )?.name ?? ""}
          onChange={(e) => {
            const preset = CUSTOM_PRESETS.find((p) => p.name === e.target.value);
            if (preset) loadPreset(preset);
          }}
        >
          <option value="" disabled>presets...</option>
          {CUSTOM_PRESETS.map((p) => (
            <option key={p.name} value={p.name}>{p.name}</option>
          ))}
        </select>
        <span style={{ fontSize: 10, color: "#667" }}>
          vars: zr zi cr ci re im
        </span>
      </div>
      <div style={{ display: "grid", gridTemplateColumns: "max-content 1fr max-content 1fr", gap: "3px 6px", alignItems: "center" }}>
        <span style={{ fontSize: 10, color: "#a0a0b0" }}>z0.re</span>
        <input
          style={exprInputStyle}
          value={def.initZr}
          onChange={(e) => updateField("initZr", e.target.value)}
        />
        <span style={{ fontSize: 10, color: "#a0a0b0" }}>z0.im</span>
        <input
          style={exprInputStyle}
          value={def.initZi}
          onChange={(e) => updateField("initZi", e.target.value)}
        />
        <span style={{ fontSize: 10, color: "#a0a0b0" }}>step.re</span>
        <input
          style={exprInputStyle}
          value={def.stepZr}
          onChange={(e) => updateField("stepZr", e.target.value)}
        />
        <span style={{ fontSize: 10, color: "#a0a0b0" }}>step.im</span>
        <input
          style={exprInputStyle}
          value={def.stepZi}
          onChange={(e) => updateField("stepZi", e.target.value)}
        />
        <span style={{ fontSize: 10, color: "#a0a0b0" }}>escape</span>
        <input
          style={{ ...exprInputStyle, gridColumn: "2 / -1" }}
          value={def.escapeExpr}
          onChange={(e) => updateField("escapeExpr", e.target.value)}
        />
      </div>
      {error && (
        <div style={{ fontSize: 10, color: "#f88", fontFamily: "ui-monospace, monospace", whiteSpace: "pre-wrap" }}>
          {error}
        </div>
      )}
    </div>
  );
};

// ----------------------------------------------------------------------------
// Julia Animation Panel
// ----------------------------------------------------------------------------
const AnimationPanel: React.FC<{
  state: ExplorerState;
  addKeyframe: () => void;
  removeKeyframe: (i: number) => void;
  updateDuration: (i: number, d: number) => void;
  isAnimating: boolean;
  onPlay: () => void;
  onStop: () => void;
  recording: boolean;
  toggleRecording: () => void;
}> = ({ state, addKeyframe, removeKeyframe, updateDuration, isAnimating, onPlay, onStop, recording, toggleRecording }) => {
  const disabled = state.fractal !== "julia";
  return (
    <div style={panelStyle}>
      <div style={{ fontSize: 11, color: "#aab", minWidth: 130 }}>
        Julia animation
        {disabled && <div style={{ color: "#776", fontSize: 10 }}>Switch to Julia to use</div>}
      </div>
      <button style={buttonStyle} onClick={addKeyframe} disabled={disabled}>
        + Keyframe
      </button>
      <div style={{ display: "flex", gap: 4, flexWrap: "wrap", maxHeight: 64, overflowY: "auto" }}>
        {state.keyframes.map((k, i) => (
          <div
            key={i}
            style={{
              background: "#1a1a25",
              border: "1px solid #2c2c3a",
              borderRadius: 4,
              padding: "2px 6px",
              display: "flex",
              alignItems: "center",
              gap: 4,
              fontSize: 10,
              fontFamily: "ui-monospace, monospace",
            }}
          >
            <span>
              #{i}: {k.c.re.toFixed(3)}, {k.c.im.toFixed(3)}
            </span>
            <input
              type="number"
              step={0.1}
              min={0.1}
              max={10}
              value={k.duration}
              style={{ ...inputStyle, width: 50 }}
              onChange={(e) => updateDuration(i, parseFloat(e.target.value) || 1)}
            />
            <span>s</span>
            <button
              onClick={() => removeKeyframe(i)}
              style={{ ...buttonStyle, padding: "1px 6px", background: "#442" }}
            >
              ×
            </button>
          </div>
        ))}
      </div>
      <button
        style={buttonStyle}
        onClick={isAnimating ? onStop : onPlay}
        disabled={disabled || state.keyframes.length < 2}
      >
        {isAnimating ? "Stop" : "Play animation"}
      </button>
      <button
        style={buttonStyle}
        onClick={toggleRecording}
        disabled={disabled || state.keyframes.length < 2}
      >
        {recording ? "● Recording (click to arm)" : "Arm record"}
      </button>
      {recording && <span style={{ color: "#f8a", fontSize: 11 }}>record armed</span>}
    </div>
  );
};

// ----------------------------------------------------------------------------
// Music Sequencer Panel
// ----------------------------------------------------------------------------
const ROOT_NAMES = ["C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"];
function rootMidiToLabel(midi: number): string {
  const pc = ((midi - 12) % 12 + 12) % 12;
  const oct = Math.floor((midi - 12) / 12);
  return `${ROOT_NAMES[pc]}${oct}`;
}

const MusicPanel: React.FC<{
  state: ExplorerState;
  setState: React.Dispatch<React.SetStateAction<ExplorerState>>;
  notes: GeneratedNote[];
  playingIdx: number;
  isPlaying: boolean;
  onGenerate: () => void;
  onPlay: () => void;
  onStop: () => void;
}> = ({ state, setState, notes, playingIdx, isPlaying, onGenerate, onPlay, onStop }) => {
  const pitches = notes.map((n) => n.pitch);
  const minPitch = pitches.length ? Math.min(...pitches) - 2 : 48;
  const maxPitch = pitches.length ? Math.max(...pitches) + 2 : 72;
  const pitchRange = maxPitch - minPitch;
  const rollHeight = 80;
  const totalBeats = notes.reduce((a, n) => a + n.duration, 0);

  return (
    <div style={{ ...panelStyle, flexDirection: "column", alignItems: "stretch" }}>
      <div style={{ display: "flex", gap: 12, flexWrap: "wrap", alignItems: "center" }}>
        <div style={{ fontSize: 11, color: "#aab" }}>Music sequencer</div>
        <label style={labelStyle}>
          scale
          <select
            style={inputStyle}
            value={state.music.scale}
            onChange={(e) =>
              setState({ ...state, music: { ...state.music, scale: e.target.value as ScaleKind } })
            }
          >
            {SCALE_OPTIONS.map((s) => (
              <option key={s} value={s}>
                {s.replace("_", " ")}
              </option>
            ))}
          </select>
        </label>
        <label style={labelStyle}>
          root: {rootMidiToLabel(state.music.rootMidi)}
          <input
            type="range"
            min={36}
            max={84}
            value={state.music.rootMidi}
            onChange={(e) =>
              setState({ ...state, music: { ...state.music, rootMidi: parseInt(e.target.value) } })
            }
          />
        </label>
        <label style={labelStyle}>
          sequence
          <select
            style={inputStyle}
            value={state.music.sequenceType}
            onChange={(e) =>
              setState({
                ...state,
                music: { ...state.music, sequenceType: e.target.value as SequenceType },
              })
            }
          >
            <option value="path">path (pitches)</option>
            <option value="region">region (chords)</option>
          </select>
        </label>
        <label style={labelStyle}>
          length: {state.music.length}
          <input
            type="range"
            min={8}
            max={64}
            value={state.music.length}
            onChange={(e) =>
              setState({ ...state, music: { ...state.music, length: parseInt(e.target.value) } })
            }
          />
        </label>
        <label style={labelStyle}>
          waveform
          <select
            style={inputStyle}
            value={state.music.waveform}
            onChange={(e) =>
              setState({ ...state, music: { ...state.music, waveform: e.target.value as Waveform } })
            }
          >
            {WAVEFORM_OPTIONS.map((w) => (
              <option key={w} value={w}>
                {w}
              </option>
            ))}
          </select>
        </label>
        <button style={buttonStyle} onClick={onGenerate}>
          Generate
        </button>
        <button style={buttonStyle} onClick={isPlaying ? onStop : onPlay} disabled={notes.length === 0 && !isPlaying}>
          {isPlaying ? "Stop" : "Play"}
        </button>
      </div>

      <div
        style={{
          height: rollHeight,
          background: "#06060c",
          borderRadius: 4,
          position: "relative",
          overflow: "hidden",
          border: "1px solid #22222c",
        }}
      >
        {notes.length === 0 && (
          <div
            style={{
              position: "absolute",
              inset: 0,
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              color: "#556",
              fontSize: 11,
            }}
          >
            click Generate to sample notes from the current fractal region
          </div>
        )}
        {(() => {
          let accBeats = 0;
          return notes.map((n, i) => {
            const x = totalBeats > 0 ? (accBeats / totalBeats) * 100 : 0;
            const w = totalBeats > 0 ? (n.duration / totalBeats) * 100 : 0;
            const normPitch = pitchRange > 0 ? (n.pitch - minPitch) / pitchRange : 0.5;
            const y = (1 - normPitch) * (rollHeight - 8) + 4;
            accBeats += n.duration;
            const active = i === playingIdx;
            return (
              <div
                key={i}
                title={`${n.name} v${n.velocity}`}
                style={{
                  position: "absolute",
                  left: `${x}%`,
                  top: y,
                  width: `${Math.max(0.5, w)}%`,
                  height: 6,
                  background: active ? "#ffd86a" : "#6adcff",
                  boxShadow: active ? "0 0 6px #ffd86a" : undefined,
                  borderRadius: 2,
                  opacity: 0.4 + (n.velocity / 127) * 0.6,
                }}
              />
            );
          });
        })()}
        {isPlaying && playingIdx >= 0 && totalBeats > 0 && (
          <div
            style={{
              position: "absolute",
              top: 0,
              bottom: 0,
              width: 1,
              background: "#ffd86a",
              left: `${(notes.slice(0, playingIdx + 1).reduce((a, n) => a + n.duration, 0) / totalBeats) * 100}%`,
              pointerEvents: "none",
            }}
          />
        )}
      </div>
    </div>
  );
};

export default FractalExplorer;
