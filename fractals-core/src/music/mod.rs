//! Music theory module: pitch classes, intervals, scales, chords, time
//! signatures, sequences — plus fractal-derived sequence generators.
//!
//! Ported from `go/src/forms/music/theory` with the same nomenclature; the
//! fractal-specific sequences (`path_sequence`, `region_sequence`) are new.

pub mod chord;
pub mod interval;
pub mod pitch;
pub mod scale;
pub mod sequence;
pub mod signature;

pub use chord::{Chord, ChordQuality};
pub use interval::Interval;
pub use pitch::{Note, NoteValue, PitchClass, Tone};
pub use scale::{Scale, ScaleKind};
pub use sequence::{
    ChordCycle, ChordSequence, FractalPathSequence, FractalRegionChordSequence, NoteSequence,
    PathPoint, RegionStats, SampledNote,
};
pub use signature::TimeSignature;
