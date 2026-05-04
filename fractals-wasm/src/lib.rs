//! WASM bridge exposing fractals-core to JavaScript.
//!
//! Design notes:
//! - `render_fractal` renders a full region and returns RGBA pixel bytes ready
//!   to be written into a canvas `ImageData`. Colours come from a palette
//!   selected by name; palettes are interpolated in OKLAB on the JS side to
//!   avoid pulling colour-science into the WASM binary (and duplicating the
//!   logic that has to exist in JS anyway for the UI previews). WASM emits the
//!   normalized smooth-iteration count and a palette key; JS applies the
//!   perceptual interpolation before drawing. See `packages/explorer/src/color.ts`.
//! - `render_fractal_progressive` renders at a lower resolution by sampling
//!   every `stride`th pixel and replicating the value into an NxN block. This
//!   is how the "1/8 -> 1/4 -> 1/2 -> full" refinement passes are built.
//! - `generate_music` runs fractal-derived `NoteSequence` / `ChordSequence`
//!   generators and returns `{ pitch, duration, velocity }` objects.
//! - Multi-threading is delegated to a JavaScript Worker pool that shards
//!   horizontal tile strips across workers, each of which owns its own WASM
//!   instance. We expose `render_tile_strip` as the building block.

use num_complex::Complex64;
use serde::{Deserialize, Serialize};
use wasm_bindgen::prelude::*;

use fractals_core::fractal::{
    BurningShip, ComplexFractal, Fractal, FractalKind, Julia, Mandelbrot, Newton,
};
use fractals_core::music::{
    FractalPathSequence, FractalRegionChordSequence, Scale, ScaleKind,
};
use fractals_core::space::{Bounds, Resolution};

#[wasm_bindgen(start)]
pub fn _start() {
    // Best-effort panic hook; pure Rust, no extra crates needed.
    std::panic::set_hook(Box::new(|info| {
        let msg = info.to_string();
        web_log(&msg);
    }));
}

#[wasm_bindgen(inline_js = r#"
export function __fractals_log(s) {
    try { console.error(s); } catch (e) {}
}
"#)]
extern "C" {
    #[wasm_bindgen(js_name = __fractals_log)]
    fn web_log(s: &str);
}

// ------ Rendering --------------------------------------------------------

fn build_fractal(kind: &str, c_re: f64, c_im: f64) -> Fractal {
    Fractal::from_name(kind, c_re, c_im)
}

/// Write smooth iteration data into a Float32Array layout: [smooth, escaped, 0, 0] per px.
/// JS colourizes using the palette. The reason we keep two slots for escaped/root is so
/// Newton's fractal can encode basin index (0..2) for distinct root colouring.
fn write_smooth_buf(
    f: &dyn ComplexFractal,
    bounds: Bounds,
    resolution: Resolution,
    max_iter: u32,
    out: &mut Vec<f32>,
) {
    let r2 = f.escape_radius();
    for py in 0..resolution.height {
        for px in 0..resolution.width {
            let c = bounds.pixel_to_complex(px, py, resolution);
            let r = f.iterate(c, max_iter);
            let smooth = r.smooth(max_iter, r2) / max_iter as f64;
            out.push(smooth as f32);
            // magnitude doubles as "basin index" for Newton; JS can use it.
            out.push(r.magnitude as f32);
        }
    }
}

/// Render the full region at full resolution. Returns `Float32Array` of length
/// width*height*2: [smooth_t, extra, smooth_t, extra, ...].
#[wasm_bindgen]
pub fn render_fractal(
    kind: &str,
    c_re: f64,
    c_im: f64,
    x_min: f64,
    x_max: f64,
    y_min: f64,
    y_max: f64,
    width: u32,
    height: u32,
    max_iter: u32,
) -> Vec<f32> {
    let f = build_fractal(kind, c_re, c_im);
    let bounds = Bounds::new(x_min, x_max, y_min, y_max);
    let resolution = Resolution::new(width, height);
    let mut out = Vec::with_capacity((width * height * 2) as usize);
    write_smooth_buf(f.as_trait(), bounds, resolution, max_iter, &mut out);
    out
}

/// Render by sampling every `stride`th pixel and replicating into NxN blocks.
/// Produces an output buffer of the same full width*height*2 size so the
/// caller can drop it straight into an ImageData of the target resolution.
#[wasm_bindgen]
pub fn render_fractal_progressive(
    kind: &str,
    c_re: f64,
    c_im: f64,
    x_min: f64,
    x_max: f64,
    y_min: f64,
    y_max: f64,
    width: u32,
    height: u32,
    max_iter: u32,
    stride: u32,
) -> Vec<f32> {
    let f = build_fractal(kind, c_re, c_im);
    let fractal = f.as_trait();
    let bounds = Bounds::new(x_min, x_max, y_min, y_max);
    let resolution = Resolution::new(width, height);
    let stride = stride.max(1);

    let sampled_w = (width + stride - 1) / stride;
    let sampled_h = (height + stride - 1) / stride;

    // First pass: compute a coarse grid.
    let mut coarse: Vec<f32> = Vec::with_capacity((sampled_w * sampled_h * 2) as usize);
    let r2 = fractal.escape_radius();
    for sy in 0..sampled_h {
        for sx in 0..sampled_w {
            let px = (sx * stride).min(width - 1);
            let py = (sy * stride).min(height - 1);
            let c = bounds.pixel_to_complex(px, py, resolution);
            let r = fractal.iterate(c, max_iter);
            let smooth = r.smooth(max_iter, r2) / max_iter as f64;
            coarse.push(smooth as f32);
            coarse.push(r.magnitude as f32);
        }
    }

    // Second pass: replicate coarse into the full buffer (nearest-neighbour).
    let mut out = vec![0f32; (width * height * 2) as usize];
    for py in 0..height {
        let sy = (py / stride).min(sampled_h - 1);
        for px in 0..width {
            let sx = (px / stride).min(sampled_w - 1);
            let src = ((sy * sampled_w + sx) * 2) as usize;
            let dst = ((py * width + px) * 2) as usize;
            out[dst] = coarse[src];
            out[dst + 1] = coarse[src + 1];
        }
    }
    out
}

/// Render a horizontal strip: rows [y_start, y_end). The caller is expected to
/// run this in parallel across multiple Workers (each with its own WASM
/// instance) to saturate CPU cores. The output buffer covers that strip only.
#[wasm_bindgen]
pub fn render_tile_strip(
    kind: &str,
    c_re: f64,
    c_im: f64,
    x_min: f64,
    x_max: f64,
    y_min: f64,
    y_max: f64,
    width: u32,
    height: u32,
    y_start: u32,
    y_end: u32,
    max_iter: u32,
) -> Vec<f32> {
    let f = build_fractal(kind, c_re, c_im);
    let fractal = f.as_trait();
    let bounds = Bounds::new(x_min, x_max, y_min, y_max);
    let resolution = Resolution::new(width, height);
    let y_end = y_end.min(height);
    let y_start = y_start.min(y_end);
    let rows = (y_end - y_start) as usize;
    let mut out = Vec::with_capacity(rows * width as usize * 2);
    let r2 = fractal.escape_radius();
    for py in y_start..y_end {
        for px in 0..width {
            let c = bounds.pixel_to_complex(px, py, resolution);
            let r = fractal.iterate(c, max_iter);
            let smooth = r.smooth(max_iter, r2) / max_iter as f64;
            out.push(smooth as f32);
            out.push(r.magnitude as f32);
        }
    }
    out
}

#[derive(Serialize)]
pub struct PointInfo {
    pub re: f64,
    pub im: f64,
    pub iterations: u32,
    pub escaped: bool,
    pub smooth: f64,
}

/// Hover inspection: compute iteration data for a single complex coordinate.
#[wasm_bindgen]
pub fn compute_point(
    kind: &str,
    c_re: f64,
    c_im: f64,
    re: f64,
    im: f64,
    max_iter: u32,
) -> JsValue {
    let f = build_fractal(kind, c_re, c_im);
    let fr = f.as_trait();
    let r = fr.iterate(Complex64::new(re, im), max_iter);
    let info = PointInfo {
        re,
        im,
        iterations: r.iterations,
        escaped: r.escaped,
        smooth: r.smooth(max_iter, fr.escape_radius()) / max_iter as f64,
    };
    serde_wasm_bindgen::to_value(&info).unwrap_or(JsValue::NULL)
}

// ------ Default bounds ----------------------------------------------------

#[derive(Serialize)]
pub struct BoundsOut {
    pub x_min: f64,
    pub x_max: f64,
    pub y_min: f64,
    pub y_max: f64,
}

#[wasm_bindgen]
pub fn default_bounds(kind: &str) -> JsValue {
    let f = build_fractal(kind, 0.0, 0.0);
    let b = match f.0 {
        FractalKind::Mandelbrot => Mandelbrot.default_bounds(),
        FractalKind::Julia(_) => Julia::default().default_bounds(),
        FractalKind::BurningShip => BurningShip.default_bounds(),
        FractalKind::Newton => Newton.default_bounds(),
    };
    let out = BoundsOut { x_min: b.x_min, x_max: b.x_max, y_min: b.y_min, y_max: b.y_max };
    serde_wasm_bindgen::to_value(&out).unwrap_or(JsValue::NULL)
}

// ------ Music generation --------------------------------------------------

#[derive(Serialize)]
pub struct NoteOut {
    pub pitch: u8,
    /// Duration in beats (quarter-note = 1).
    pub duration: f64,
    pub velocity: u8,
    /// Human-readable tone name ("C4", etc) for UI display.
    pub name: String,
}

#[derive(Deserialize)]
pub struct MusicConfig {
    pub kind: String,
    pub c_re: f64,
    pub c_im: f64,
    pub scale: String,
    pub root_midi: u8,
    pub sequence_type: String, // "path" | "region"
    pub length: u32,
    pub max_iter: u32,
    pub x_min: f64,
    pub x_max: f64,
    pub y_min: f64,
    pub y_max: f64,
}

fn dispatch_music(cfg: &MusicConfig) -> Vec<NoteOut> {
    let scale_kind = ScaleKind::parse(&cfg.scale);
    let root_tone = fractals_core::music::pitch::Tone::from_midi(cfg.root_midi);
    let scale = Scale::new(root_tone, scale_kind);
    let bounds = Bounds::new(cfg.x_min, cfg.x_max, cfg.y_min, cfg.y_max);
    let length = cfg.length.max(1) as usize;

    match cfg.sequence_type.as_str() {
        "region" => {
            let build_chords = |f: &dyn ComplexFractal| -> Vec<NoteOut> {
                // Arpeggiate chord progression as individual notes for playback.
                let regions = 8.min(length);
                let regions = regions.max(2);
                let mut notes = Vec::with_capacity(length);
                // We need to dispatch on concrete type for FractalRegionChordSequence,
                // so match kind at construction.
                match cfg.kind.as_str() {
                    "julia" => {
                        let mut seq = FractalRegionChordSequence::new(
                            Julia::new(Complex64::new(cfg.c_re, cfg.c_im)),
                            bounds,
                            regions,
                            scale.clone(),
                            cfg.max_iter,
                        );
                        build_region_notes(&mut seq, &mut notes, length);
                    }
                    "burning_ship" => {
                        let mut seq = FractalRegionChordSequence::new(
                            BurningShip,
                            bounds,
                            regions,
                            scale.clone(),
                            cfg.max_iter,
                        );
                        build_region_notes(&mut seq, &mut notes, length);
                    }
                    "newton" => {
                        let mut seq = FractalRegionChordSequence::new(
                            Newton,
                            bounds,
                            regions,
                            scale.clone(),
                            cfg.max_iter,
                        );
                        build_region_notes(&mut seq, &mut notes, length);
                    }
                    _ => {
                        let mut seq = FractalRegionChordSequence::new(
                            Mandelbrot,
                            bounds,
                            regions,
                            scale.clone(),
                            cfg.max_iter,
                        );
                        build_region_notes(&mut seq, &mut notes, length);
                    }
                }
                // silence unused warning in non-matching branch
                let _ = f;
                notes
            };
            let f = Fractal::from_name(&cfg.kind, cfg.c_re, cfg.c_im);
            build_chords(f.as_trait())
        }
        _ => {
            // Default: path-based note sequence.
            let start = Complex64::new(cfg.x_min, (cfg.y_min + cfg.y_max) * 0.5);
            let end = Complex64::new(cfg.x_max, (cfg.y_min + cfg.y_max) * 0.5);
            let mut notes = Vec::with_capacity(length);
            match cfg.kind.as_str() {
                "julia" => {
                    let seq = FractalPathSequence::new(
                        Julia::new(Complex64::new(cfg.c_re, cfg.c_im)),
                        start,
                        end,
                        length,
                        scale,
                        cfg.max_iter,
                    );
                    for i in 0..length {
                        let s = seq.sample_at(i);
                        notes.push(NoteOut {
                            pitch: s.tone.midi_code(),
                            duration: 4.0 / s.value as f64,
                            velocity: s.velocity,
                            name: s.tone.as_string(),
                        });
                    }
                }
                "burning_ship" => {
                    let seq = FractalPathSequence::new(
                        BurningShip,
                        start,
                        end,
                        length,
                        scale,
                        cfg.max_iter,
                    );
                    for i in 0..length {
                        let s = seq.sample_at(i);
                        notes.push(NoteOut {
                            pitch: s.tone.midi_code(),
                            duration: 4.0 / s.value as f64,
                            velocity: s.velocity,
                            name: s.tone.as_string(),
                        });
                    }
                }
                "newton" => {
                    let seq =
                        FractalPathSequence::new(Newton, start, end, length, scale, cfg.max_iter);
                    for i in 0..length {
                        let s = seq.sample_at(i);
                        notes.push(NoteOut {
                            pitch: s.tone.midi_code(),
                            duration: 4.0 / s.value as f64,
                            velocity: s.velocity,
                            name: s.tone.as_string(),
                        });
                    }
                }
                _ => {
                    let seq = FractalPathSequence::new(
                        Mandelbrot,
                        start,
                        end,
                        length,
                        scale,
                        cfg.max_iter,
                    );
                    for i in 0..length {
                        let s = seq.sample_at(i);
                        notes.push(NoteOut {
                            pitch: s.tone.midi_code(),
                            duration: 4.0 / s.value as f64,
                            velocity: s.velocity,
                            name: s.tone.as_string(),
                        });
                    }
                }
            }
            notes
        }
    }
}

fn build_region_notes<F: ComplexFractal>(
    seq: &mut FractalRegionChordSequence<F>,
    notes: &mut Vec<NoteOut>,
    length: usize,
) {
    // Cycle through the region chord progression as arpeggiated notes.
    use fractals_core::music::ChordSequence;
    let mut i: usize = 0;
    while notes.len() < length {
        let c = seq.next_chord();
        for tone in &c.tones {
            if notes.len() >= length {
                break;
            }
            // Varying velocity per chord keeps playback lively.
            let velocity = 70u8 + ((i as u8) % 4) * 10;
            notes.push(NoteOut {
                pitch: tone.midi_code(),
                duration: 0.5,
                velocity,
                name: tone.as_string(),
            });
        }
        i += 1;
    }
}

#[wasm_bindgen]
pub fn generate_music(config: JsValue) -> JsValue {
    let cfg: MusicConfig = match serde_wasm_bindgen::from_value(config) {
        Ok(c) => c,
        Err(e) => {
            web_log(&format!("generate_music config parse failed: {e}"));
            return JsValue::NULL;
        }
    };
    let notes = dispatch_music(&cfg);
    serde_wasm_bindgen::to_value(&notes).unwrap_or(JsValue::NULL)
}

// Expose a string list of supported fractal kinds for the UI.
#[wasm_bindgen]
pub fn fractal_kinds() -> JsValue {
    let kinds = vec!["mandelbrot", "julia", "burning_ship", "newton"];
    serde_wasm_bindgen::to_value(&kinds).unwrap_or(JsValue::NULL)
}
