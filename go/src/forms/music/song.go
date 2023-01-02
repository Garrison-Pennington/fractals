package music

import (
	mt "fractals/forms/music/theory"

	"gitlab.com/gomidi/midi/gm"
	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

type Instrument struct {
	Name string
	GM   int
}

type Song struct {
	Tempo float64
	Meter mt.TimeSignature
	Clock smf.MetricTicks
	SMF   *smf.SMF
}

func (s Song) SetTempo(tempo float64) Song {
	s.Tempo = tempo
	return s
}

func (s Song) SetMeter(meter mt.TimeSignature) Song {
	s.Meter = meter
	return s
}

func (s Song) PlayQuarters(notes []uint8) smf.Track {
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(s.Tempo))
	tr.Add(0, smf.MetaMeter(s.Meter.Beats, s.Meter.Value))
	tr.Add(0, smf.MetaInstrument("Piano"))
	tr.Add(0, midi.ProgramChange(0, gm.Instr_AcousticGrandPiano.Value()))
	for i := range notes {
		tr = Play4th(tr, s.Clock, 0, notes[i], 100)
	}
	return tr
}

func (s Song) AddTrack(track smf.Track) {
	track.Close(0)
	s.SMF.Add(track)
}

func (s Song) Save(filename string) {
	s.SMF.WriteFile(filename)
}

func BasicSong() Song {
	return Song{
		Tempo: 140,
		Meter: mt.TimeSignature{4, 4},
		Clock: smf.MetricTicks(960),
		SMF:   smf.New(),
	}
}

func (s Song) PianoTrack() smf.Track {
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(s.Tempo))
	tr.Add(0, smf.MetaMeter(s.Meter.Beats, s.Meter.Value))
	tr.Add(0, smf.MetaInstrument("Piano"))
	tr.Add(0, midi.ProgramChange(0, gm.Instr_AcousticGrandPiano.Value()))
	return tr
}
