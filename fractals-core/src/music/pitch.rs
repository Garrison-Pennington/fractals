//! Pitch classes, tones, notes.

/// A pitch class (e.g. "C", "F#") with its MIDI base (the MIDI code in octave 0
/// — following the Go conventions; A is 21, C is 12).
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct PitchClass {
    pub name: &'static str,
    pub midi_base: u8,
}

impl PitchClass {
    pub const fn new(name: &'static str, midi_base: u8) -> Self {
        Self { name, midi_base }
    }
    pub fn tone(self, octave: u8) -> Tone {
        Tone { pitch_class: self, octave }
    }
}

// Pitch class constants — matching the Go base values.
pub const A: PitchClass = PitchClass::new("A", 21);
pub const A_SHARP: PitchClass = PitchClass::new("A#", 22);
pub const B: PitchClass = PitchClass::new("B", 23);
pub const C: PitchClass = PitchClass::new("C", 12);
pub const C_SHARP: PitchClass = PitchClass::new("C#", 13);
pub const D: PitchClass = PitchClass::new("D", 14);
pub const D_SHARP: PitchClass = PitchClass::new("D#", 15);
pub const E: PitchClass = PitchClass::new("E", 16);
pub const F: PitchClass = PitchClass::new("F", 17);
pub const F_SHARP: PitchClass = PitchClass::new("F#", 18);
pub const G: PitchClass = PitchClass::new("G", 19);
pub const G_SHARP: PitchClass = PitchClass::new("G#", 20);

pub const PITCH_CLASSES: [PitchClass; 12] =
    [C, C_SHARP, D, D_SHARP, E, F, F_SHARP, G, G_SHARP, A, A_SHARP, B];

/// A specific pitch — pitch class + octave.
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
pub struct Tone {
    pub pitch_class: PitchClass,
    pub octave: u8,
}

impl Tone {
    /// Standard MIDI pitch code (C0 = 12, A4 = 69, etc).
    pub fn midi_code(&self) -> u8 {
        self.pitch_class.midi_base.saturating_add(12u8.saturating_mul(self.octave))
    }

    pub fn from_midi(code: u8) -> Self {
        // C-relative octave numbering; mirrors the Go implementation.
        let pc_idx = ((code as i32) - 12).rem_euclid(12) as usize;
        let octave = ((code as i32 - 12) / 12).max(0) as u8;
        Self { pitch_class: PITCH_CLASSES[pc_idx], octave }
    }

    pub fn pitched_up(&self, half_steps: u8) -> Self {
        Self::from_midi(self.midi_code().saturating_add(half_steps))
    }
    pub fn pitched_down(&self, half_steps: u8) -> Self {
        Self::from_midi(self.midi_code().saturating_sub(half_steps))
    }
    pub fn octave_up(&self, octaves: u8) -> Self {
        Self { pitch_class: self.pitch_class, octave: self.octave.saturating_add(octaves) }
    }
    pub fn octave_down(&self, octaves: u8) -> Self {
        Self { pitch_class: self.pitch_class, octave: self.octave.saturating_sub(octaves) }
    }
    pub fn as_string(&self) -> String {
        format!("{}{}", self.pitch_class.name, self.octave)
    }
    /// Frequency in Hz using A4 = 440.
    pub fn frequency(&self) -> f64 {
        440.0 * 2f64.powf((self.midi_code() as f64 - 69.0) / 12.0)
    }
}

/// Duration denominator (1 = whole, 2 = half, 4 = quarter, ...).
pub type NoteValue = u8;

#[derive(Debug, Clone, Copy)]
pub struct Note {
    pub tone: Tone,
    pub value: NoteValue,
}

impl Note {
    pub fn new(tone: Tone, value: NoteValue) -> Self {
        Self { tone, value }
    }
    pub fn whole(tone: Tone) -> Self {
        Self::new(tone, 1)
    }
    pub fn half(tone: Tone) -> Self {
        Self::new(tone, 2)
    }
    pub fn quarter(tone: Tone) -> Self {
        Self::new(tone, 4)
    }
    pub fn eighth(tone: Tone) -> Self {
        Self::new(tone, 8)
    }
    /// Duration in beats where a quarter note is 1 beat.
    pub fn duration_beats(&self) -> f64 {
        4.0 / self.value as f64
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn midi_round_trip() {
        let t = Tone { pitch_class: C, octave: 4 };
        assert_eq!(t.midi_code(), 60);
        let back = Tone::from_midi(60);
        assert_eq!(back.pitch_class, C);
        assert_eq!(back.octave, 4);
    }

    #[test]
    fn pitched_up_crosses_octave() {
        let t = Tone { pitch_class: B, octave: 3 };
        let up = t.pitched_up(2);
        assert_eq!(up.pitch_class, C_SHARP);
        assert_eq!(up.octave, 4);
    }

    #[test]
    fn frequency_a4() {
        let t = Tone { pitch_class: A, octave: 4 };
        assert!((t.frequency() - 440.0).abs() < 1e-9);
    }
}
