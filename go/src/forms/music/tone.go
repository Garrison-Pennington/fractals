package music

import "strconv"

type Tone struct {
	PitchClass PitchClass
	Octave     uint8
}

// CONSTRUCTORS:
func FromMidiCode(code uint8) Tone {
	tone := MIDI_MODS[(code-A.MidiBase)%12]
	octave := (code - C.MidiBase) / 12
	return Tone{tone, octave}
}

// METHODS:
func (n Tone) PitchedUp(halfSteps uint8) Tone {
	return FromMidiCode(n.MidiCode() + halfSteps)
}

func (n Tone) PitchedDown(halfSteps uint8) Tone {
	return FromMidiCode(n.MidiCode() - halfSteps)
}

func (n Tone) OctaveUp(octaves uint8) Tone {
	return Tone{n.PitchClass, n.Octave + octaves}
}

func (n Tone) OctaveDown(octaves uint8) Tone {
	return Tone{n.PitchClass, n.Octave - octaves}
}

func (n Tone) Distance(next Tone) uint8 {
	return n.MidiCode() - next.MidiCode()
}

func (n Tone) MidiCode() uint8 {
	return n.PitchClass.MidiBase + 12*n.Octave
}

func (n Tone) AsString() string {
	return n.PitchClass.Value + strconv.FormatInt(int64(n.Octave), 10)
}

func (n Tone) Major() Chord {
	return MajorTriad(n)
}

func (n Tone) Minor() Chord {
	return MinorTriad(n)
}

func ListTones(notes []Tone) (str string) {
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
func (t PitchClass) Tone(octave uint8) Tone {
	return Tone{t, octave}
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
