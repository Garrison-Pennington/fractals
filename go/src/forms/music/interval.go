package music

type Interval struct {
	HalfSteps uint8
}

func (i Interval) NextTone(note Tone) (next Tone) {
	return note.PitchedUp(i.HalfSteps)
}

var MINOR_2ND = Interval{1}
var MAJOR_2ND = Interval{2}
var MINOR_3RD = Interval{3}
var MAJOR_3RD = Interval{4}
var PERFECT_4TH = Interval{5}
var DIMINISHED_5TH = Interval{6}
var PERFECT_5TH = Interval{7}
var AUGMENTED_5TH = Interval{8}
var MINOR_6TH = Interval{8}
var MAJOR_6TH = Interval{9}
var DIMINISHED_7TH = Interval{9}
var MINOR_7TH = Interval{10}
var MAJOR_7TH = Interval{11}
