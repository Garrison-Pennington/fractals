//! Time signatures.

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct TimeSignature {
    pub beats: u8,
    pub value: u8,
}

impl TimeSignature {
    pub const fn new(beats: u8, value: u8) -> Self {
        Self { beats, value }
    }
    pub fn ratio(&self) -> f64 {
        self.beats as f64 / self.value as f64
    }
}

pub const COMMON_TIME: TimeSignature = TimeSignature::new(4, 4);
pub const WALTZ_TIME: TimeSignature = TimeSignature::new(3, 4);

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn common_time_is_one() {
        assert!((COMMON_TIME.ratio() - 1.0).abs() < 1e-12);
    }
}
