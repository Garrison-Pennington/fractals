//! Benchmark render throughput (megapixels/sec) for each fractal at 1000 iters.

use std::time::Instant;

use criterion::{criterion_group, criterion_main, Criterion};

use fractals_core::fractal::{render, BurningShip, ComplexFractal, Julia, Mandelbrot, Newton};
use fractals_core::space::{Bounds, Resolution};

const MAX_ITER: u32 = 1000;
const W: u32 = 256;
const H: u32 = 256;

fn bench_fractal<F: ComplexFractal>(c: &mut Criterion, name: &str, f: F, bounds: Bounds) {
    let res = Resolution::new(W, H);
    let total_px = (W as u64) * (H as u64);
    c.bench_function(&format!("render_{}_{}x{}_{}iter", name, W, H, MAX_ITER), |b| {
        b.iter_custom(|iters| {
            let start = Instant::now();
            for _ in 0..iters {
                let _ = render(&f, bounds, res, MAX_ITER);
            }
            let elapsed = start.elapsed();
            let mpx_per_sec =
                (iters as f64 * total_px as f64) / elapsed.as_secs_f64() / 1_000_000.0;
            println!("  {} throughput: {:.3} Mpx/s", name, mpx_per_sec);
            elapsed
        });
    });
}

fn benches(c: &mut Criterion) {
    bench_fractal(c, "mandelbrot", Mandelbrot, Mandelbrot.default_bounds());
    let j = Julia::default();
    let jb = j.default_bounds();
    bench_fractal(c, "julia", j, jb);
    bench_fractal(c, "burning_ship", BurningShip, BurningShip.default_bounds());
    bench_fractal(c, "newton", Newton, Newton.default_bounds());
}

criterion_group!(render_bench, benches);
criterion_main!(render_bench);
