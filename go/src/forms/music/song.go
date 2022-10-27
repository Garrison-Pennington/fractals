package music

import (
	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

type Song struct {
	Tempo  float64
	Meter  TimeSignature
	Clock  smf.MetricTicks
	SMF    *smf.SMF
	Tracks []smf.Track
	fresh  bool
}

func MakeSong(tempo float64, meter TimeSignature, clock smf.MetricTicks) Song {
	var tr smf.Track
	tr.Add(0, smf.MetaTempo(tempo))
	tr.Add(0, smf.MetaMeter(meter.Beats, meter.Value))
	tracks := []smf.Track{tr}
	return Song{
		Tempo:  tempo,
		Meter:  meter,
		Clock:  clock,
		SMF:    smf.New(),
		Tracks: tracks,
		fresh:  true,
	}
}

func (s *Song) newTrack(instrument Instrument) *smf.Track {
	if s.fresh {
		s.Tracks[0].Add(0, smf.MetaInstrument(instrument.Name))
		s.Tracks[0].Add(0, midi.ProgramChange(0, instrument.GM))
		s.fresh = false
	} else {
		var track smf.Track
		track.Add(0, smf.MetaInstrument(instrument.Name))
		track.Add(0, midi.ProgramChange(uint8(len(s.Tracks)), instrument.GM))
		s.Tracks = append(s.Tracks, track)
	}
	return &s.Tracks[len(s.Tracks)-1]
}

func (s Song) CurrentChannel() uint8 {
	return uint8(len(s.Tracks) - 1)
}

func (s *Song) PlayQuarters(notes []uint8, instrument Instrument) {
	tr := s.newTrack(instrument)
	for i := range notes {
		Play4th(tr, s.Clock, s.CurrentChannel(), notes[i], 100)
	}
}

func (s *Song) PlayWholes(notes []uint8, instrument Instrument) {
	tr := s.newTrack(instrument)
	for i := range notes {
		PlayWhole(tr, s.Clock, s.CurrentChannel(), notes[i], 100)
	}
}

func (s *Song) PlayWithLengths(notes []uint8, lengths []uint16, instrument Instrument) {
	tr := s.newTrack(instrument)
	PlayNotesWithLengths(notes, lengths, tr, s.Clock, s.CurrentChannel(), 100)
}

func (s Song) Save(filename string) {
	for _, tr := range s.Tracks {
		tr.Close(0)
		s.SMF.Add(tr)
	}
	s.SMF.WriteFile(filename)
}

func BasicSong() Song {
	return MakeSong(140, TimeSignature{4, 4}, smf.MetricTicks(960))
}
