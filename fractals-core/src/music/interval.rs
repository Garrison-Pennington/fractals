//! Intervals measured in half-steps.

use super::pitch::Tone;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Interval(pub u8);

impl Interval {
    pub fn next_tone(&self, t: Tone) -> Tone {
        t.pitched_up(self.0)
    }
    pub fn last_tone(&self, t: Tone) -> Tone {
        t.pitched_down(self.0)
    }
}

pub const MINOR_2ND: Interval = Interval(1);
pub const MAJOR_2ND: Interval = Interval(2);
pub const MINOR_3RD: Interval = Interval(3);
pub const MAJOR_3RD: Interval = Interval(4);
pub const PERFECT_4TH: Interval = Interval(5);
pub const DIMINISHED_5TH: Interval = Interval(6);
pub const PERFECT_5TH: Interval = Interval(7);
pub const AUGMENTED_5TH: Interval = Interval(8);
pub const MINOR_6TH: Interval = Interval(8);
pub const MAJOR_6TH: Interval = Interval(9);
pub const DIMINISHED_7TH: Interval = Interval(9);
pub const MINOR_7TH: Interval = Interval(10);
pub const MAJOR_7TH: Interval = Interval(11);
pub const OCTAVE: Interval = Interval(12);

#[cfg(test)]
mod tests {
    use super::super::pitch::{C, E};
    use super::*;

    #[test]
    fn major_third_up_from_c4_is_e4() {
        let c4 = C.tone(4);
        let e4 = MAJOR_3RD.next_tone(c4);
        assert_eq!(e4.pitch_class, E);
        assert_eq!(e4.octave, 4);
    }
}
