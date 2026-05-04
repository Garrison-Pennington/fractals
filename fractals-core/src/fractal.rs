//! Complex-plane fractals using escape-time iteration.
//!
//! `ComplexFractal` is the trait every fractal implements. `iterate` takes a
//! complex coordinate plus a max-iteration cap and returns both the iteration
//! count and the escape magnitude, which together drive smooth coloring.

use num_complex::Complex64;

use crate::space::{Bounds, Resolution};

/// Result of iterating a fractal from a starting complex coordinate.
///
/// - `iterations`: how many steps ran before the orbit escaped (or `max` if it
///   never did).
/// - `magnitude`: |z| at the final iteration — used for smooth coloring.
/// - `escaped`: true iff the orbit left the escape radius before `max`.
#[derive(Debug, Clone, Copy, PartialEq)]
pub struct IterationResult {
    pub iterations: u32,
    pub magnitude: f64,
    pub escaped: bool,
}

impl IterationResult {
    /// Normalized iteration count with log smoothing — the classic smooth
    /// coloring formula. Returns a value in `[0, max_iter as f64]` that is
    /// continuous across iteration bands.
    pub fn smooth(&self, max_iter: u32, escape_radius: f64) -> f64 {
        if !self.escaped || self.iterations >= max_iter {
            return max_iter as f64;
        }
        let log_zn = self.magnitude.max(1.0).ln();
        let nu = (log_zn / escape_radius.ln()).ln() / 2f64.ln();
        (self.iterations as f64 + 1.0 - nu).clamp(0.0, max_iter as f64)
    }
}

/// Core trait: every fractal implements `iterate`.
pub trait ComplexFractal: Sync {
    fn name(&self) -> &'static str;
    fn escape_radius(&self) -> f64 {
        2.0
    }
    fn iterate(&self, c: Complex64, max_iter: u32) -> IterationResult;
    fn default_bounds(&self) -> Bounds {
        Bounds::new(-2.0, 2.0, -2.0, 2.0)
    }
}

/// Classic Mandelbrot set: z_{n+1} = z_n^2 + c, z_0 = 0.
#[derive(Debug, Clone, Copy, Default)]
pub struct Mandelbrot;

impl ComplexFractal for Mandelbrot {
    fn name(&self) -> &'static str {
        "Mandelbrot"
    }
    fn default_bounds(&self) -> Bounds {
        Bounds::new(-2.5, 1.0, -1.25, 1.25)
    }
    fn iterate(&self, c: Complex64, max_iter: u32) -> IterationResult {
        let mut z = Complex64::new(0.0, 0.0);
        let r2 = self.escape_radius().powi(2);
        for i in 0..max_iter {
            if z.norm_sqr() > r2 {
                return IterationResult { iterations: i, magnitude: z.norm(), escaped: true };
            }
            z = z * z + c;
        }
        IterationResult { iterations: max_iter, magnitude: z.norm(), escaped: false }
    }
}

/// Julia set for z_{n+1} = z_n^2 + c, where c is fixed and z_0 = pixel.
#[derive(Debug, Clone, Copy)]
pub struct Julia {
    pub c: Complex64,
}

impl Julia {
    pub fn new(c: Complex64) -> Self {
        Self { c }
    }
}

impl Default for Julia {
    fn default() -> Self {
        Self { c: Complex64::new(-0.7, 0.27015) }
    }
}

impl ComplexFractal for Julia {
    fn name(&self) -> &'static str {
        "Julia"
    }
    fn default_bounds(&self) -> Bounds {
        Bounds::new(-1.5, 1.5, -1.5, 1.5)
    }
    fn iterate(&self, c: Complex64, max_iter: u32) -> IterationResult {
        let mut z = c;
        let r2 = self.escape_radius().powi(2);
        for i in 0..max_iter {
            if z.norm_sqr() > r2 {
                return IterationResult { iterations: i, magnitude: z.norm(), escaped: true };
            }
            z = z * z + self.c;
        }
        IterationResult { iterations: max_iter, magnitude: z.norm(), escaped: false }
    }
}

/// Burning Ship: z_{n+1} = (|Re z_n| + i|Im z_n|)^2 + c.
#[derive(Debug, Clone, Copy, Default)]
pub struct BurningShip;

impl ComplexFractal for BurningShip {
    fn name(&self) -> &'static str {
        "BurningShip"
    }
    fn default_bounds(&self) -> Bounds {
        Bounds::new(-2.2, 1.5, -2.0, 1.2)
    }
    fn iterate(&self, c: Complex64, max_iter: u32) -> IterationResult {
        let mut z = Complex64::new(0.0, 0.0);
        let r2 = self.escape_radius().powi(2);
        for i in 0..max_iter {
            if z.norm_sqr() > r2 {
                return IterationResult { iterations: i, magnitude: z.norm(), escaped: true };
            }
            let zr = z.re.abs();
            let zi = z.im.abs();
            let folded = Complex64::new(zr, zi);
            z = folded * folded + c;
        }
        IterationResult { iterations: max_iter, magnitude: z.norm(), escaped: false }
    }
}

/// Newton's fractal for z^3 - 1 = 0.
///
/// Convergence (rather than escape) drives the iteration count; `magnitude`
/// encodes which root was found (0, 1, 2 → angle bucket) so the WASM layer can
/// colour each basin distinctly.
#[derive(Debug, Clone, Copy, Default)]
pub struct Newton;

impl Newton {
    const TOLERANCE: f64 = 1e-6;
    fn roots() -> [Complex64; 3] {
        let t = std::f64::consts::TAU / 3.0;
        [
            Complex64::new(1.0, 0.0),
            Complex64::new(t.cos(), t.sin()),
            Complex64::new((2.0 * t).cos(), (2.0 * t).sin()),
        ]
    }
}

impl ComplexFractal for Newton {
    fn name(&self) -> &'static str {
        "Newton"
    }
    fn escape_radius(&self) -> f64 {
        // Not really an escape radius; reused for smooth coloring.
        2.0
    }
    fn default_bounds(&self) -> Bounds {
        Bounds::new(-2.0, 2.0, -2.0, 2.0)
    }
    fn iterate(&self, c: Complex64, max_iter: u32) -> IterationResult {
        let roots = Self::roots();
        let mut z = c;
        if z.norm() < 1e-10 {
            // Derivative is zero; treat as not-converged.
            return IterationResult { iterations: max_iter, magnitude: 0.0, escaped: false };
        }
        for i in 0..max_iter {
            let z2 = z * z;
            let num = z2 * z - Complex64::new(1.0, 0.0);
            let den = Complex64::new(3.0, 0.0) * z2;
            if den.norm() < 1e-12 {
                return IterationResult {
                    iterations: max_iter,
                    magnitude: 0.0,
                    escaped: false,
                };
            }
            z -= num / den;
            for (idx, root) in roots.iter().enumerate() {
                if (z - *root).norm() < Self::TOLERANCE {
                    return IterationResult {
                        iterations: i,
                        // Encode root index in magnitude (0.0, 1.0, 2.0) so callers
                        // know which basin this point landed in.
                        magnitude: idx as f64,
                        escaped: true,
                    };
                }
            }
        }
        IterationResult { iterations: max_iter, magnitude: 3.0, escaped: false }
    }
}

/// Enum wrapper that lets callers pick a fractal by tag. The WASM layer uses
/// this; library-only callers can still use the concrete types directly.
#[derive(Debug, Clone, Copy)]
pub enum FractalKind {
    Mandelbrot,
    Julia(Julia),
    BurningShip,
    Newton,
}

/// Trait-object-ish newtype so we can dispatch based on kind without boxing.
pub struct Fractal(pub FractalKind);

impl Fractal {
    pub fn from_name(name: &str, c_re: f64, c_im: f64) -> Self {
        match name {
            "julia" => Fractal(FractalKind::Julia(Julia::new(Complex64::new(c_re, c_im)))),
            "burning_ship" | "burningship" => Fractal(FractalKind::BurningShip),
            "newton" => Fractal(FractalKind::Newton),
            _ => Fractal(FractalKind::Mandelbrot),
        }
    }

    pub fn as_trait(&self) -> &dyn ComplexFractal {
        match &self.0 {
            FractalKind::Mandelbrot => &Mandelbrot,
            FractalKind::Julia(j) => j,
            FractalKind::BurningShip => &BurningShip,
            FractalKind::Newton => &Newton,
        }
    }
}

/// Run iteration for every pixel in `resolution` within `bounds`.
/// Returns a flat row-major Vec of `IterationResult`.
pub fn render(
    f: &dyn ComplexFractal,
    bounds: Bounds,
    resolution: Resolution,
    max_iter: u32,
) -> Vec<IterationResult> {
    let mut out = Vec::with_capacity(resolution.total_pixels() as usize);
    for py in 0..resolution.height {
        for px in 0..resolution.width {
            let c = bounds.pixel_to_complex(px, py, resolution);
            out.push(f.iterate(c, max_iter));
        }
    }
    out
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn mandelbrot_interior_never_escapes() {
        // (0, 0) is inside the Mandelbrot set.
        let r = Mandelbrot.iterate(Complex64::new(0.0, 0.0), 1000);
        assert_eq!(r.iterations, 1000);
        assert!(!r.escaped);
    }

    #[test]
    fn mandelbrot_exterior_escapes_fast() {
        // (2, 2) escapes on the first step.
        let r = Mandelbrot.iterate(Complex64::new(2.0, 2.0), 1000);
        assert!(r.escaped);
        assert!(r.iterations < 5);
    }

    #[test]
    fn julia_default_produces_varied_output() {
        let j = Julia::default();
        let r1 = j.iterate(Complex64::new(0.0, 0.0), 500);
        let r2 = j.iterate(Complex64::new(2.0, 2.0), 500);
        assert!(r2.escaped);
        assert_ne!(r1.iterations, r2.iterations);
    }

    #[test]
    fn burning_ship_far_point_escapes() {
        let r = BurningShip.iterate(Complex64::new(3.0, 3.0), 1000);
        assert!(r.escaped);
    }

    #[test]
    fn newton_converges_to_a_root() {
        // Near the real root (1, 0), should converge in a few iterations.
        let r = Newton.iterate(Complex64::new(1.1, 0.0), 500);
        assert!(r.escaped);
        assert_eq!(r.magnitude, 0.0); // First root index.
    }

    #[test]
    fn smooth_value_is_monotonic_for_varying_orbit() {
        let r1 = Mandelbrot.iterate(Complex64::new(-0.5, 0.0), 100);
        let r2 = Mandelbrot.iterate(Complex64::new(-2.5, 0.0), 100);
        // r1 stays inside; r2 escapes immediately.
        let s1 = r1.smooth(100, 2.0);
        let s2 = r2.smooth(100, 2.0);
        assert!(s1 > s2);
    }

    #[test]
    fn render_produces_expected_dimensions() {
        let out = render(
            &Mandelbrot,
            Bounds::new(-2.0, 1.0, -1.0, 1.0),
            Resolution::new(10, 8),
            100,
        );
        assert_eq!(out.len(), 80);
    }

    #[test]
    fn render_non_uniform() {
        let out = render(
            &Mandelbrot,
            Mandelbrot.default_bounds(),
            Resolution::new(100, 100),
            100,
        );
        let first = out[0].iterations;
        assert!(out.iter().any(|r| r.iterations != first));
    }
}
