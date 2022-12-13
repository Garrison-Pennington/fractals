package music

type Measure struct {
	Signature TimeSignature
	Notes     []struct {
		Note
		uint32
	}
	Resolution uint32
}

type Note struct {
	Tone  Tone
	Value uint8
}

func NewMeasure(sig TimeSignature, resolution uint32) Measure {
	return Measure{sig, make([]struct {
		Note
		uint32
	}, 0), resolution}
}

func (m *Measure) AddNote(note Note, at uint32) {
	m.Notes = append(m.Notes, struct {
		Note
		uint32
	}{note, at})
}
