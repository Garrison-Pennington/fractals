package music

import (
	mt "fractals/forms/music/theory"

	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

func AddNotes(tr smf.Track, clock smf.MetricTicks, channel uint8, velocity uint8, notes ...mt.Note) smf.Track {
	ons, offs := mt.MidiMessages(notes, channel, velocity)
	for _, on := range ons {
		tr.Add(0, on)
	}
	delay := (clock.Ticks4th() * 4) / uint32(notes[0].Value)
	for _, off := range offs {
		tr.Add(delay, off)
		delay = 0
	}
	return tr
}

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
