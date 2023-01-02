package music

type ToneSequence interface {
	Next() Tone
}

func NextTones(seq ToneSequence, n int) (tones []Tone) {
	for i := 0; i < n; i++ {
		tones = append(tones, seq.Next())
	}
	return
}
