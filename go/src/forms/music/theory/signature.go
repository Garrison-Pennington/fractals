package music_theory

import smf "gitlab.com/gomidi/midi/v2/smf"

type TimeSignature struct {
	Beats uint8
	Value uint8
}

func (ts TimeSignature) Whole(clock smf.MetricTicks) uint32 {
	return clock.Ticks4th() * 4
}

func (ts TimeSignature) Half(clock smf.MetricTicks) uint32 {
	return clock.Ticks4th() * 2
}

func (ts TimeSignature) Quarter(clock smf.MetricTicks) uint32 {
	return clock.Ticks4th()
}

func (ts TimeSignature) Eigth(clock smf.MetricTicks) uint32 {
	return clock.Ticks8th()
}

func (ts TimeSignature) Sixteenth(clock smf.MetricTicks) uint32 {
	return clock.Ticks16th()
}

func (ts TimeSignature) Ratio() float64 {
	return float64(ts.Beats) / float64(ts.Value)
}
