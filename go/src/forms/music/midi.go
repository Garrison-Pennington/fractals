package music

import (
	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

func PlayNote(channel uint8, note uint8, velocity uint8) (midi.Message, midi.Message) {
	return midi.NoteOn(channel, note, velocity), midi.NoteOff(channel, note)
}

func PlayNotes(channel uint8, notes []uint8, velocity uint8) (ons []midi.Message, offs []midi.Message) {
	for _, note := range notes {
		on, off := PlayNote(channel, note, velocity)
		ons = append(ons, on)
		offs = append(offs, off)
	}
	return
}

func Play16th(tr smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) smf.Track {
	on, off := PlayNote(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(clock.Ticks16th(), off)
	return tr
}

func Play8th(tr smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) smf.Track {
	on, off := PlayNote(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(clock.Ticks8th(), off)
	return tr
}

func Play4th(tr smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) smf.Track {
	on, off := PlayNote(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(clock.Ticks4th(), off)
	return tr
}

func PlayHalf(tr smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) smf.Track {
	on, off := PlayNote(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(clock.Ticks4th()*2, off)
	return tr
}

func PlayWhole(tr smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) smf.Track {
	on, off := PlayNote(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(clock.Ticks4th()*4, off)
	return tr
}
