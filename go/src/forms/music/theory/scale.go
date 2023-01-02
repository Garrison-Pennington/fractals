package music

type Scale struct {
	Root      Tone
	Intervals []Interval
}

// CONSTRUCTORS:
func MinorScale(root Tone) Scale {
	return Scale{root, MINOR_SCALE_INTERVALS}
}

func MajorScale(root Tone) Scale {
	return Scale{root, MAJOR_SCALE_INTERVALS}
}

// METHODS:
func (s Scale) I() Tone {
	return s.Root
}

func (s Scale) II() Tone {
	return s.Intervals[0].NextTone(s.Root)
}

func (s Scale) III() Tone {
	return s.Intervals[1].NextTone(s.Root)
}

func (s Scale) IV() Tone {
	return s.Intervals[2].NextTone(s.Root)
}

func (s Scale) V() Tone {
	return s.Intervals[3].NextTone(s.Root)
}

func (s Scale) VI() Tone {
	return s.Intervals[4].NextTone(s.Root)
}

func (s Scale) VII() Tone {
	return s.Intervals[5].NextTone(s.Root)
}

func (s Scale) VIII() Tone {
	return s.Intervals[6].NextTone(s.Root)
}

func (s Scale) Tones(sequence []int) (tones []Tone) {
	for _, val := range sequence {
		var tone Tone
		if mod := val % 7; mod < 0 {
			tone = s.Intervals[6+mod].LastTone(s.Root)
		} else {
			tone = s.Intervals[mod].NextTone(s.Root)
		}
		if octaves := val / 7; octaves < 0 {
			tone = tone.OctaveDown(uint8(octaves * -1))
		} else {
			tone = tone.OctaveUp(uint8(octaves))
		}
		tones = append(tones, tone)
	}
	return
}

// CONSTANTS:
var MINOR_SCALE_INTERVALS []Interval = []Interval{MINOR_2ND, MINOR_3RD, PERFECT_4TH, PERFECT_5TH, MINOR_6TH, MINOR_7TH, OCTAVE}
var MAJOR_SCALE_INTERVALS []Interval = []Interval{MAJOR_2ND, MAJOR_3RD, PERFECT_4TH, PERFECT_5TH, MAJOR_6TH, MAJOR_7TH, OCTAVE}
