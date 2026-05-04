//! Chords — triads, sevenths, extended.

use super::interval::*;
use super::pitch::{PitchClass, Tone};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ChordQuality {
    Major,
    Minor,
    Augmented,
    Diminished,
    DominantSeventh,
    MajorSeventh,
    MinorSeventh,
    DiminishedSeventh,
    HalfDiminishedSeventh,
    MajorNinth,
    MinorNinth,
}

impl ChordQuality {
    pub fn suffix(&self) -> &'static str {
        match self {
            ChordQuality::Major => "",
            ChordQuality::Minor => "m",
            ChordQuality::Augmented => "+",
            ChordQuality::Diminished => "°",
            ChordQuality::DominantSeventh => "7",
            ChordQuality::MajorSeventh => "maj7",
            ChordQuality::MinorSeventh => "m7",
            ChordQuality::DiminishedSeventh => "°7",
            ChordQuality::HalfDiminishedSeventh => "ø7",
            ChordQuality::MajorNinth => "maj9",
            ChordQuality::MinorNinth => "m9",
        }
    }
}

#[derive(Debug, Clone)]
pub struct Chord {
    pub root: PitchClass,
    pub tones: Vec<Tone>,
    pub quality: ChordQuality,
}

impl Chord {
    pub fn new(root: Tone, quality: ChordQuality) -> Self {
        let tones = match quality {
            ChordQuality::Major => {
                vec![root, MAJOR_3RD.next_tone(root), PERFECT_5TH.next_tone(root)]
            }
            ChordQuality::Minor => {
                vec![root, MINOR_3RD.next_tone(root), PERFECT_5TH.next_tone(root)]
            }
            ChordQuality::Augmented => {
                vec![root, MAJOR_3RD.next_tone(root), AUGMENTED_5TH.next_tone(root)]
            }
            ChordQuality::Diminished => {
                vec![root, MINOR_3RD.next_tone(root), DIMINISHED_5TH.next_tone(root)]
            }
            ChordQuality::DominantSeventh => vec![
                root,
                MAJOR_3RD.next_tone(root),
                PERFECT_5TH.next_tone(root),
                MINOR_7TH.next_tone(root),
            ],
            ChordQuality::MajorSeventh => vec![
                root,
                MAJOR_3RD.next_tone(root),
                PERFECT_5TH.next_tone(root),
                MAJOR_7TH.next_tone(root),
            ],
            ChordQuality::MinorSeventh => vec![
                root,
                MINOR_3RD.next_tone(root),
                PERFECT_5TH.next_tone(root),
                MINOR_7TH.next_tone(root),
            ],
            ChordQuality::DiminishedSeventh => vec![
                root,
                MINOR_3RD.next_tone(root),
                DIMINISHED_5TH.next_tone(root),
                DIMINISHED_7TH.next_tone(root),
            ],
            ChordQuality::HalfDiminishedSeventh => vec![
                root,
                MINOR_3RD.next_tone(root),
                DIMINISHED_5TH.next_tone(root),
                MINOR_7TH.next_tone(root),
            ],
            ChordQuality::MajorNinth => vec![
                root,
                MAJOR_3RD.next_tone(root),
                PERFECT_5TH.next_tone(root),
                MAJOR_7TH.next_tone(root),
                MAJOR_2ND.next_tone(root).octave_up(1),
            ],
            ChordQuality::MinorNinth => vec![
                root,
                MINOR_3RD.next_tone(root),
                PERFECT_5TH.next_tone(root),
                MINOR_7TH.next_tone(root),
                MAJOR_2ND.next_tone(root).octave_up(1),
            ],
        };
        Self { root: root.pitch_class, tones, quality }
    }

    pub fn major_triad(root: Tone) -> Self {
        Self::new(root, ChordQuality::Major)
    }
    pub fn minor_triad(root: Tone) -> Self {
        Self::new(root, ChordQuality::Minor)
    }

    pub fn name(&self) -> String {
        format!("{}{}", self.root.name, self.quality.suffix())
    }

    /// Rotate bass note up an octave (first inversion, in place).
    pub fn invert(&mut self) {
        if self.tones.is_empty() {
            return;
        }
        let first = self.tones.remove(0);
        self.tones.push(first.octave_up(1));
    }
}

#[cfg(test)]
mod tests {
    use super::super::pitch::{C, E, G};
    use super::*;

    #[test]
    fn c_major_triad_is_c_e_g() {
        let c = Chord::major_triad(C.tone(4));
        assert_eq!(c.tones[0].pitch_class, C);
        assert_eq!(c.tones[1].pitch_class, E);
        assert_eq!(c.tones[2].pitch_class, G);
    }

    #[test]
    fn c_dom7_has_bb() {
        let c = Chord::new(C.tone(4), ChordQuality::DominantSeventh);
        assert_eq!(c.tones.len(), 4);
        assert_eq!(c.tones[3].midi_code(), C.tone(4).midi_code() + 10);
    }

    #[test]
    fn inversion_rotates() {
        let mut c = Chord::major_triad(C.tone(4));
        let before = c.tones[0];
        c.invert();
        assert_eq!(c.tones.last().unwrap().pitch_class, before.pitch_class);
    }
}
