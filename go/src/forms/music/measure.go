package music

type Measure struct {
	Signature TimeSignature
	Notes     map[uint8][]Note
}

type Note struct {
	Tone  Tone
	Value uint8
}

func NewMeasure(sig TimeSignature) Measure {
	return Measure{sig, make(map[uint8][]Note)}
}

func (m *Measure) AddNotes(at uint8, notes ...Note) {
	m.Notes[at] = append(m.Notes[at], notes...)
}

func (m *Measure) AddNoteSequence(notes []Note, at uint8) {
	for _, note := range notes {
		m.AddNotes(at, note)
		at += noteDuration(note)
	}
}

func (m Measure) onBeats() (count uint8) {
	offFlags := registerOffFlags(m.Notes)
	var notesOn uint8 = 0
	var i uint8 = 0
	for i < 32 {
		notesOn += countNotes(m.Notes, i)
		notesOn -= countFlags(offFlags, i)
		if notesOn > 0 {
			count++
		}
		i++
	}
	return
}

func registerOffFlags(noteMap map[uint8][]Note) (flags map[uint8]uint8) {
	for start, notes := range noteMap {
		for _, note := range notes {
			flags[start+noteDuration(note)] += 1
		}
	}
	return
}

func countNotes(noteMap map[uint8][]Note, i uint8) (count uint8) {
	if val, ok := noteMap[i]; ok {
		count += uint8(len(val))
	}
	return
}

func countFlags(noteMap map[uint8]uint8, i uint8) (count uint8) {
	if val, ok := noteMap[i]; ok {
		count += uint8(val)
	}
	return
}

func noteDuration(note Note) uint8 {
	return uint8(32 / note.Value)
}
