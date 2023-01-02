package music

import (
	mt "fractals/forms/music/theory"
)

type ToneSequence interface {
	Next() mt.Tone
}

func NextTones(seq ToneSequence, n int) (tones []mt.Tone) {
	for i := 0; i < n; i++ {
		tones = append(tones, seq.Next())
	}
	return
}

type ChordSequence interface {
	Next() mt.Chord
}

func NextChords(seq ChordSequence, n int) (chords []mt.Chord) {
	for i := 0; i < n; i++ {
		chords = append(chords, seq.Next())
	}
	return
}
