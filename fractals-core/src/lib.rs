//! fractals-core: Fractal computation engine and fractal-derived music sequences.
//!
//! Ported from the original Go implementation in `go/src/forms/` with added
//! fractal variants (Julia, Burning Ship, Newton) and music sequences that
//! derive pitches / chord progressions from fractal iteration data.

pub mod fractal;
pub mod music;
pub mod space;

pub use fractal::{
    BurningShip, ComplexFractal, Fractal, FractalKind, IterationResult, Julia, Mandelbrot, Newton,
};
pub use space::{Bounds, Resolution};
