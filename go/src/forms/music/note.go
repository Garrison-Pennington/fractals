package music

import "strconv"

type Note struct {
	PitchClass PitchClass
	Octave     uint8
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
	return Note{n.PitchClass, n.Octave + octaves}
}

func (n Note) OctaveDown(octaves uint8) Note {
	return Note{n.PitchClass, n.Octave - octaves}
}

func (n Note) Distance(next Note) uint8 {
	return n.MidiCode() - next.MidiCode()
}

func (n Note) MidiCode() uint8 {
	return n.PitchClass.MidiBase + 12*n.Octave
}

func (n Note) AsString() string {
	return n.PitchClass.Value + strconv.FormatInt(int64(n.Octave), 10)
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

type PitchClass struct {
	Value    string
	MidiBase uint8
}

// CONSTRUCTORS:

// METHODS:
func (t PitchClass) Note(octave uint8) Note {
	return Note{t, octave}
}

// CONSTANTS:
var A PitchClass = PitchClass{"A", 21}
var AS PitchClass = PitchClass{"A#", 22}
var B PitchClass = PitchClass{"B", 23}
var C PitchClass = PitchClass{"C", 12}
var CS PitchClass = PitchClass{"C#", 13}
var D PitchClass = PitchClass{"D", 14}
var DS PitchClass = PitchClass{"D#", 15}
var E PitchClass = PitchClass{"E", 16}
var F PitchClass = PitchClass{"F", 17}
var FS PitchClass = PitchClass{"F#", 18}
var G PitchClass = PitchClass{"G", 19}
var GS PitchClass = PitchClass{"G#", 20}

var MIDI_MODS []PitchClass = []PitchClass{A, AS, B, C, CS, D, DS, E, F, FS, G, GS}
