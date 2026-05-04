//! Note / chord sequences — including fractal-derived sequence generators.
//!
//! `ChordCycle` matches the Go example. The two new sequence types derive
//! musical output directly from fractal iteration data:
//!
//! - `FractalPathSequence` samples points along a path through the complex
//!   plane, maps each sample's iteration count into a scale degree, and emits
//!   a stream of notes.
//! - `FractalRegionChordSequence` divides the visible region into tiles,
//!   averages iteration density per tile, and emits a chord progression where
//!   each chord's root + quality is chosen from that density.
//!
//! Both are deterministic given the same fractal, scale, and path/region config.

use num_complex::Complex64;

use super::chord::{Chord, ChordQuality};
use super::pitch::{Note, Tone};
use super::scale::Scale;
use crate::fractal::ComplexFractal;
use crate::space::Bounds;

pub trait NoteSequence {
    fn next_note(&mut self) -> Note;
    fn take(&mut self, n: usize) -> Vec<Note> {
        (0..n).map(|_| self.next_note()).collect()
    }
}

pub trait ChordSequence {
    fn next_chord(&mut self) -> Chord;
    fn take(&mut self, n: usize) -> Vec<Chord> {
        (0..n).map(|_| self.next_chord()).collect()
    }
}

/// Cycle through a list of chords with each chord repeated N times.
#[derive(Debug, Clone)]
pub struct ChordCycle {
    pub chords: Vec<Chord>,
    pub repetitions: u8,
    rep: u8,
    idx: usize,
}

impl ChordCycle {
    pub fn new(chords: Vec<Chord>, repetitions: u8) -> Self {
        Self { chords, repetitions: repetitions.max(1), rep: 0, idx: 0 }
    }
}

impl ChordSequence for ChordCycle {
    fn next_chord(&mut self) -> Chord {
        let chord = self.chords[self.idx].clone();
        self.rep += 1;
        if self.rep >= self.repetitions {
            self.rep = 0;
            self.idx = (self.idx + 1) % self.chords.len();
        }
        chord
    }
}

/// A single sampled note, returned by `FractalPathSequence::sample`.
#[derive(Debug, Clone, Copy)]
pub struct SampledNote {
    pub tone: Tone,
    pub value: u8,
    pub velocity: u8,
}

/// A point along the sampling path, with the iteration result it produced.
#[derive(Debug, Clone, Copy)]
pub struct PathPoint {
    pub c: Complex64,
    pub iterations: u32,
}

/// Maps fractal iteration counts sampled along a straight-line path into scale
/// degrees.
///
/// Construct with `new(fractal, path_start, path_end, length, scale, max_iter)`.
/// Each `next_note()` call advances along the path (uniform spacing) and emits
/// one note whose pitch is determined by the sampled iteration count. The
/// iterator wraps back to the start after `length` samples so long sequences
/// repeat deterministically.
pub struct FractalPathSequence<F: ComplexFractal> {
    pub fractal: F,
    pub start: Complex64,
    pub end: Complex64,
    pub length: usize,
    pub scale: Scale,
    pub max_iter: u32,
    pub base_velocity: u8,
    pub note_value: u8,
    cursor: usize,
}

impl<F: ComplexFractal> FractalPathSequence<F> {
    pub fn new(
        fractal: F,
        start: Complex64,
        end: Complex64,
        length: usize,
        scale: Scale,
        max_iter: u32,
    ) -> Self {
        Self {
            fractal,
            start,
            end,
            length: length.max(1),
            scale,
            max_iter,
            base_velocity: 80,
            note_value: 8,
            cursor: 0,
        }
    }

    pub fn sample_at(&self, i: usize) -> SampledNote {
        let t = if self.length <= 1 { 0.0 } else { i as f64 / (self.length - 1) as f64 };
        let c = Complex64::new(
            self.start.re + (self.end.re - self.start.re) * t,
            self.start.im + (self.end.im - self.start.im) * t,
        );
        let r = self.fractal.iterate(c, self.max_iter);
        // Map iteration count onto [0, scale_len * octave_span).
        let scale_len = self.scale.len() as i32;
        let span = scale_len * 3; // 3-octave range
        let degree = if self.max_iter == 0 {
            0
        } else {
            ((r.iterations as i64 * span as i64) / self.max_iter as i64) as i32 % span
        };
        let tone = self.scale.degree(degree);
        // Louder when orbit stayed close to the set (high iteration count).
        let velocity = ((r.iterations as u32 * 127) / self.max_iter.max(1)) as u8;
        let velocity = velocity.max(30);
        SampledNote { tone, value: self.note_value, velocity }
    }

    pub fn path_points(&self) -> Vec<PathPoint> {
        (0..self.length)
            .map(|i| {
                let t = if self.length <= 1 { 0.0 } else { i as f64 / (self.length - 1) as f64 };
                let c = Complex64::new(
                    self.start.re + (self.end.re - self.start.re) * t,
                    self.start.im + (self.end.im - self.start.im) * t,
                );
                let r = self.fractal.iterate(c, self.max_iter);
                PathPoint { c, iterations: r.iterations }
            })
            .collect()
    }
}

impl<F: ComplexFractal> NoteSequence for FractalPathSequence<F> {
    fn next_note(&mut self) -> Note {
        let sampled = self.sample_at(self.cursor);
        self.cursor = (self.cursor + 1) % self.length;
        Note::new(sampled.tone, sampled.value)
    }
}

/// Statistics for one region of the complex plane.
#[derive(Debug, Clone, Copy)]
pub struct RegionStats {
    pub avg_iterations: f64,
    pub escape_fraction: f64,
}

/// Derives a chord progression from the average iteration density of regions
/// sampled across `bounds`.
///
/// The visible region is divided into `regions` tiles arranged horizontally.
/// For each tile, we average iteration counts and use the normalized mean to
/// pick a scale degree (chord root) and a quality (major / minor / 7th) from a
/// small lookup — producing a coherent, data-driven progression.
pub struct FractalRegionChordSequence<F: ComplexFractal> {
    pub fractal: F,
    pub bounds: Bounds,
    pub regions: usize,
    pub scale: Scale,
    pub max_iter: u32,
    pub samples_per_region: u32,
    pub root_octave: u8,
    cursor: usize,
}

impl<F: ComplexFractal> FractalRegionChordSequence<F> {
    pub fn new(
        fractal: F,
        bounds: Bounds,
        regions: usize,
        scale: Scale,
        max_iter: u32,
    ) -> Self {
        Self {
            fractal,
            bounds,
            regions: regions.max(1),
            scale,
            max_iter,
            samples_per_region: 64,
            root_octave: 4,
            cursor: 0,
        }
    }

    pub fn region_stats(&self, region_idx: usize) -> RegionStats {
        let regions = self.regions as f64;
        let region_w = self.bounds.width() / regions;
        let x_start = self.bounds.x_min + region_idx as f64 * region_w;
        let x_end = x_start + region_w;
        let samples = self.samples_per_region;
        // Fixed 2D grid sampling — deterministic for a given config.
        let side = (samples as f64).sqrt().max(2.0) as u32;
        let mut total_iter: u64 = 0;
        let mut escapes: u32 = 0;
        let mut count: u32 = 0;
        for sy in 0..side {
            let v = (sy as f64 + 0.5) / side as f64;
            let y = self.bounds.y_min + v * self.bounds.height();
            for sx in 0..side {
                let u = (sx as f64 + 0.5) / side as f64;
                let x = x_start + u * (x_end - x_start);
                let r = self.fractal.iterate(Complex64::new(x, y), self.max_iter);
                total_iter += r.iterations as u64;
                if r.escaped {
                    escapes += 1;
                }
                count += 1;
            }
        }
        let count = count.max(1) as f64;
        RegionStats {
            avg_iterations: total_iter as f64 / count,
            escape_fraction: escapes as f64 / count,
        }
    }

    fn chord_for_region(&self, stats: RegionStats) -> Chord {
        let norm = (stats.avg_iterations / self.max_iter as f64).clamp(0.0, 1.0);
        let scale_len = self.scale.len() as i32;
        let degree = (norm * (scale_len - 1) as f64).round() as i32;
        let root = self.scale.degree(degree);
        // More escapes => brighter (major / dominant); fewer => darker (minor / 7ths).
        let quality = match (stats.escape_fraction * 5.0) as u32 {
            0 => ChordQuality::MinorSeventh,
            1 => ChordQuality::Minor,
            2 => ChordQuality::Major,
            3 => ChordQuality::DominantSeventh,
            _ => ChordQuality::MajorSeventh,
        };
        Chord::new(Tone { pitch_class: root.pitch_class, octave: self.root_octave }, quality)
    }
}

impl<F: ComplexFractal> ChordSequence for FractalRegionChordSequence<F> {
    fn next_chord(&mut self) -> Chord {
        let stats = self.region_stats(self.cursor);
        let chord = self.chord_for_region(stats);
        self.cursor = (self.cursor + 1) % self.regions;
        chord
    }
}

#[cfg(test)]
mod tests {
    use super::super::pitch::C;
    use super::*;
    use crate::fractal::Mandelbrot;
    use crate::music::scale::ScaleKind;

    #[test]
    fn chord_cycle_repeats_then_advances() {
        let c = Chord::major_triad(C.tone(4));
        let mut cyc = ChordCycle::new(vec![c.clone(), c.clone()], 2);
        assert_eq!(cyc.next_chord().name(), "C");
        assert_eq!(cyc.next_chord().name(), "C");
        // After two reps on the 1st entry we advance to the 2nd (still C in this fixture).
        assert_eq!(cyc.next_chord().name(), "C");
    }

    #[test]
    fn path_sequence_is_deterministic() {
        let scale = Scale::new(C.tone(4), ScaleKind::Major);
        let mut a = FractalPathSequence::new(
            Mandelbrot,
            Complex64::new(-2.0, 0.0),
            Complex64::new(1.0, 0.0),
            16,
            scale.clone(),
            200,
        );
        let mut b = FractalPathSequence::new(
            Mandelbrot,
            Complex64::new(-2.0, 0.0),
            Complex64::new(1.0, 0.0),
            16,
            scale,
            200,
        );
        let xs = a.take(16);
        let ys = b.take(16);
        for (x, y) in xs.iter().zip(ys.iter()) {
            assert_eq!(x.tone.midi_code(), y.tone.midi_code());
            assert_eq!(x.value, y.value);
        }
    }

    #[test]
    fn path_sequence_produces_varied_notes() {
        let scale = Scale::new(C.tone(4), ScaleKind::Major);
        let mut seq = FractalPathSequence::new(
            Mandelbrot,
            Complex64::new(-2.0, 0.1),
            Complex64::new(1.0, 0.1),
            16,
            scale,
            200,
        );
        let notes = seq.take(16);
        let unique: std::collections::HashSet<_> = notes.iter().map(|n| n.tone.midi_code()).collect();
        assert!(unique.len() > 1);
    }

    #[test]
    fn region_chord_sequence_is_deterministic() {
        let scale = Scale::new(C.tone(4), ScaleKind::Major);
        let mut a = FractalRegionChordSequence::new(
            Mandelbrot,
            Bounds::new(-2.0, 1.0, -1.0, 1.0),
            4,
            scale.clone(),
            100,
        );
        let mut b = FractalRegionChordSequence::new(
            Mandelbrot,
            Bounds::new(-2.0, 1.0, -1.0, 1.0),
            4,
            scale,
            100,
        );
        let xs = a.take(4);
        let ys = b.take(4);
        for (x, y) in xs.iter().zip(ys.iter()) {
            assert_eq!(x.name(), y.name());
        }
    }
}
