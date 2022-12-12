package music

import smf "gitlab.com/gomidi/midi/v2/smf"

type TimeSignature struct {
	Beats uint8
	Value uint8
	Clock smf.MetricTicks
}

func (ts TimeSignature) Whole() uint32 {
	return ts.Clock.Ticks4th() * 4
}

func (ts TimeSignature) Half() uint32 {
	return ts.Clock.Ticks4th() * 2
}

func (ts TimeSignature) Quarter() uint32 {
	return ts.Clock.Ticks4th()
}

func (ts TimeSignature) Eigth() uint32 {
	return ts.Clock.Ticks8th()
}

func (ts TimeSignature) Sixteenth() uint32 {
	return ts.Clock.Ticks16th()
}
