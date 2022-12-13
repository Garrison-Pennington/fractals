package music

type Measure struct {
	Signature TimeSignature
	Notes     []struct {
		Note
		uint8
	}
}

type Note struct {
	Tone  Tone
	Value uint8
}

func NewMeasure(sig TimeSignature) Measure {
	return Measure{sig, make([]struct {
		Note
		uint8
	}, 0)}
}

func (m *Measure) AddNote(note Note, at uint8) {
	m.Notes = append(m.Notes, struct {
		Note
		uint8
	}{note, at})
}
