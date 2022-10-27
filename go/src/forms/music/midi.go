package music

import (
	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
	"golang.org/x/exp/constraints"
)

func max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func PlayNoteMsgs(channel uint8, note uint8, velocity uint8) (midi.Message, midi.Message) {
	return midi.NoteOn(channel, note, velocity), midi.NoteOff(channel, note)
}

func PlayNote(tr *smf.Track, deltaTicks uint32, channel uint8, note uint8, velocity uint8) {
	on, off := PlayNoteMsgs(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(deltaTicks, off)
}

func Play16th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks16th(), channel, note, velocity)
}

func Play8th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks8th(), channel, note, velocity)
}

func Play4th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th(), channel, note, velocity)
}

func PlayHalf(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th()*2, channel, note, velocity)
}

func PlayWhole(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th()*4, channel, note, velocity)
}

func PlayNotesWithLengths(notes []uint8, lengths []uint16, tr *smf.Track, clock smf.MetricTicks, channel uint8, velocity uint8){
	ln, ll := len(notes), len(lengths)
	for i := 0; i < max(ln, ll); i++ {
		il, in := i % ll, i % ln
		NOTE_VALUES[lengths[il]](tr, clock, channel, notes[in], velocity)
	}
}
