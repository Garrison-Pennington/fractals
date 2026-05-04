//! Scales — sequences of intervals relative to a root tone.
//!
//! Supports major, minor, major pentatonic, minor pentatonic, and blues.

use super::interval::*;
use super::pitch::Tone;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ScaleKind {
    Major,
    Minor,
    MajorPentatonic,
    MinorPentatonic,
    Blues,
}

impl ScaleKind {
    pub fn intervals(&self) -> &'static [Interval] {
        match self {
            ScaleKind::Major => &MAJOR_SCALE,
            ScaleKind::Minor => &MINOR_SCALE,
            ScaleKind::MajorPentatonic => &MAJOR_PENTATONIC,
            ScaleKind::MinorPentatonic => &MINOR_PENTATONIC,
            ScaleKind::Blues => &BLUES_SCALE,
        }
    }
    pub fn parse(s: &str) -> Self {
        match s.to_ascii_lowercase().as_str() {
            "minor" => ScaleKind::Minor,
            "major_pentatonic" | "pentatonic_major" => ScaleKind::MajorPentatonic,
            "minor_pentatonic" | "pentatonic_minor" | "pentatonic" => ScaleKind::MinorPentatonic,
            "blues" => ScaleKind::Blues,
            _ => ScaleKind::Major,
        }
    }
}

#[derive(Debug, Clone)]
pub struct Scale {
    pub root: Tone,
    pub kind: ScaleKind,
}

impl Scale {
    pub fn new(root: Tone, kind: ScaleKind) -> Self {
        Self { root, kind }
    }
    pub fn major(root: Tone) -> Self {
        Self::new(root, ScaleKind::Major)
    }
    pub fn minor(root: Tone) -> Self {
        Self::new(root, ScaleKind::Minor)
    }

    /// Return the tone at a given 0-based scale degree. Wraps across octaves
    /// via negative or >=len indices, mirroring the Go `Tones` behavior.
    pub fn degree(&self, degree: i32) -> Tone {
        let intervals = self.kind.intervals();
        let len = intervals.len() as i32;
        let octaves = degree.div_euclid(len);
        let idx = degree.rem_euclid(len) as usize;
        // Degree 0 is root; degree N is root + intervals[N-1].
        let base = if idx == 0 { self.root } else { intervals[idx - 1].next_tone(self.root) };
        if octaves > 0 {
            base.octave_up(octaves as u8)
        } else if octaves < 0 {
            base.octave_down((-octaves) as u8)
        } else {
            base
        }
    }

    pub fn tones(&self, degrees: &[i32]) -> Vec<Tone> {
        degrees.iter().map(|&d| self.degree(d)).collect()
    }

    pub fn len(&self) -> usize {
        self.kind.intervals().len() + 1
    }
}

pub const MAJOR_SCALE: [Interval; 7] =
    [MAJOR_2ND, MAJOR_3RD, PERFECT_4TH, PERFECT_5TH, MAJOR_6TH, MAJOR_7TH, OCTAVE];
pub const MINOR_SCALE: [Interval; 7] =
    [MAJOR_2ND, MINOR_3RD, PERFECT_4TH, PERFECT_5TH, MINOR_6TH, MINOR_7TH, OCTAVE];
pub const MAJOR_PENTATONIC: [Interval; 5] =
    [MAJOR_2ND, MAJOR_3RD, PERFECT_5TH, MAJOR_6TH, OCTAVE];
pub const MINOR_PENTATONIC: [Interval; 5] =
    [MINOR_3RD, PERFECT_4TH, PERFECT_5TH, MINOR_7TH, OCTAVE];
pub const BLUES_SCALE: [Interval; 6] = [
    MINOR_3RD,
    PERFECT_4TH,
    DIMINISHED_5TH,
    PERFECT_5TH,
    MINOR_7TH,
    OCTAVE,
];

#[cfg(test)]
mod tests {
    use super::super::pitch::{A, C, E, G};
    use super::*;

    #[test]
    fn c_major_degrees() {
        let scale = Scale::major(C.tone(4));
        assert_eq!(scale.degree(0).pitch_class, C);
        assert_eq!(scale.degree(2).pitch_class, E);
        assert_eq!(scale.degree(4).pitch_class, G);
    }

    #[test]
    fn a_minor_has_a_root() {
        let scale = Scale::minor(A.tone(4));
        assert_eq!(scale.degree(0).pitch_class, A);
    }

    #[test]
    fn major_pentatonic_has_five_unique_pitches() {
        let scale = Scale::new(C.tone(4), ScaleKind::MajorPentatonic);
        let pcs: Vec<_> = (0..5).map(|d| scale.degree(d).pitch_class).collect();
        let mut s = pcs.clone();
        s.sort_by_key(|p| p.midi_base);
        s.dedup();
        assert_eq!(s.len(), 5);
    }

    #[test]
    fn negative_degree_goes_down_an_octave() {
        let scale = Scale::major(C.tone(4));
        let below = scale.degree(-1);
        assert!(below.midi_code() < C.tone(4).midi_code());
    }
}
