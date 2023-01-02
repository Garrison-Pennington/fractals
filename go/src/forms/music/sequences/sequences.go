package music_sequences

import (
	mt "fractals/forms/music/theory"
)

type NoteSequence interface {
	Next() mt.Note
}

func NextNotes(seq NoteSequence, n int) (notes []mt.Note) {
	for i := 0; i < n; i++ {
		notes = append(notes, seq.Next())
	}
	return
}

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
