package music

import "strconv"

type Note struct {
	Tone   Tone
	Octave uint8
}

// CONSTRUCTORS:
func FromMidiCode(code uint8) Note {
	tone := MIDI_MODS[(code-A.MidiBase)%12]
	octave := (code - C.MidiBase) / 12
	return Note{tone, octave}
}

// METHODS:
func (n Note) PitchedUp(halfSteps uint8) Note {
	return FromMidiCode(n.MidiCode() + halfSteps)
}

func (n Note) PitchedDown(halfSteps uint8) Note {
	return FromMidiCode(n.MidiCode() - halfSteps)
}

func (n Note) OctaveUp(octaves uint8) Note {
	return Note{n.Tone, n.Octave + octaves}
}

func (n Note) OctaveDown(octaves uint8) Note {
	return Note{n.Tone, n.Octave - octaves}
}

func (n Note) Distance(next Note) uint8 {
	return n.MidiCode() - next.MidiCode()
}

func (n Note) MidiCode() uint8 {
	return n.Tone.MidiBase + 12*n.Octave
}

func (n Note) AsString() string {
	return n.Tone.Value + strconv.FormatInt(int64(n.Octave), 10)
}

func (n Note) Major() Chord {
	return MajorTriad(n)
}

func (n Note) Minor() Chord {
	return MinorTriad(n)
}

func ListNotes(notes []Note) (str string) {
	for _, val := range notes {
		str += val.AsString() + " "
	}
	return
}

type Tone struct {
	Value    string
	MidiBase uint8
}

// CONSTRUCTORS:

// METHODS:
func (t Tone) Note(octave uint8) Note {
	return Note{t, octave}
}

// CONSTANTS:
var A Tone = Tone{"A", 21}
var AS Tone = Tone{"A#", 22}
var B Tone = Tone{"B", 23}
var C Tone = Tone{"C", 12}
var CS Tone = Tone{"C#", 13}
var D Tone = Tone{"D", 14}
var DS Tone = Tone{"D#", 15}
var E Tone = Tone{"E", 16}
var F Tone = Tone{"F", 17}
var FS Tone = Tone{"F#", 18}
var G Tone = Tone{"G", 19}
var GS Tone = Tone{"G#", 20}

var MIDI_MODS []Tone = []Tone{A, AS, B, C, CS, D, DS, E, F, FS, G, GS}
