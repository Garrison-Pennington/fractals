package music_theory

import midi "gitlab.com/gomidi/midi/v2"

type Chord struct {
	Root    PitchClass
	Tones   []Tone
	Quality string
}

// CONSTRUCTORS:
func MinorTriad(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root)}, "m"}
}

func MajorTriad(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root)}, ""}
}

func AugmentedTriad(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), AUGMENTED_5TH.NextTone(root)}, "+"}
}

func DiminishedTriad(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), DIMINISHED_5TH.NextTone(root)}, "o"}
}

func DiminishedSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), DIMINISHED_5TH.NextTone(root), DIMINISHED_7TH.NextTone(root)}, "o7"}
}

func HalfDiminishedSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), DIMINISHED_5TH.NextTone(root), MINOR_7TH.NextTone(root)}, root.PitchClass.Value}
}

func MinorSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root), MINOR_7TH.NextTone(root)}, "m7"}
}

func MinorMajorSeventh(root Tone, inversions uint8) Chord {
	return Chord{root.PitchClass, []Tone{root, MINOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root), MAJOR_7TH.NextTone(root)}, root.PitchClass.Value}
}

func DominantSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root), MINOR_7TH.NextTone(root)}, root.PitchClass.Value + "7"}
}

func MajorSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), PERFECT_5TH.NextTone(root), MAJOR_7TH.NextTone(root)}, root.PitchClass.Value}
}

func AugmentedSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), AUGMENTED_5TH.NextTone(root), MINOR_7TH.NextTone(root)}, root.PitchClass.Value}
}

func AugmentedMajorSeventh(root Tone) Chord {
	return Chord{root.PitchClass, []Tone{root, MAJOR_3RD.NextTone(root), AUGMENTED_5TH.NextTone(root), MAJOR_7TH.NextTone(root)}, root.PitchClass.Value}
}

// METHODS:
func (c *Chord) Invert() {
	c.Tones = append(c.Tones[1:], c.Tones[0].OctaveUp(1))
}

func (c Chord) Name() string {
	return c.Root.Value + c.Quality
}

func (c Chord) arpeggioGenerator(pattern []uint8) (ch chan Tone) {
	ch = make(chan Tone)
	go func() {
		defer close(ch)
		for i := 0; true; i++ {
			if i == len(pattern) {
				i = 0
			}
			tone := c.Tones[pattern[i]-1]
			ch <- tone
		}
	}()
	return
}

func (c Chord) Arpeggiate(pattern []uint8, count uint8) (result []Tone) {
	result = make([]Tone, count)
	var n uint8 = 0
	for tone := range c.arpeggioGenerator(pattern) {
		result[n] = tone
		n++
		if n == count {
			break
		}
	}
	return
}

func (c Chord) MidiMessages(channel uint8, velocity uint8) (ons []midi.Message, offs []midi.Message) {
	for _, tone := range c.Tones {
		on, off := tone.MidiMessages(channel, velocity)
		ons = append(ons, on)
		offs = append(offs, off)
	}
	return
}

func (c Chord) Equal(other Chord, ignoreInversions bool) bool {
	// TODO: Implement ignoreInversions
	return c.Root == other.Root && SameTones(c.Tones, other.Tones) && c.Quality == other.Quality
}

func SameChords(s1 []Chord, s2 []Chord) bool {
	for i, chord := range s1 {
		if !chord.Equal(s2[i], false) {
			return false
		}
	}
	return true
}

func ListChords(chords []Chord) (str string) {
	for _, val := range chords {
		str += val.Name() + " "
	}
	return
}

// Major Triads
var CM Chord = C.Tone(4).Major()
var CSM Chord = CS.Tone(4).Major()
var DM Chord = D.Tone(4).Major()
var DSM Chord = DS.Tone(4).Major()
var EM Chord = E.Tone(4).Major()
var FM Chord = F.Tone(4).Major()
var FSM Chord = FS.Tone(4).Major()
var GM Chord = G.Tone(4).Major()
var GSM Chord = GS.Tone(4).Major()
var AM Chord = A.Tone(4).Major()
var ASM Chord = AS.Tone(4).Major()
var BM Chord = B.Tone(4).Major()

var MAJORS []Chord = []Chord{CM, CSM, DM, DSM, EM, FM, FSM, GM, GSM, AM, ASM, BM}

// First Inversion Major Triads
// Second Inversion Major Triads
// Minor Triads
var Cm Chord = C.Tone(4).Minor()
var CSm Chord = CS.Tone(4).Minor()
var Dm Chord = D.Tone(4).Minor()
var DSm Chord = DS.Tone(4).Minor()
var Em Chord = E.Tone(4).Minor()
var Fm Chord = F.Tone(4).Minor()
var FSm Chord = FS.Tone(4).Minor()
var Gm Chord = G.Tone(4).Minor()
var GSm Chord = GS.Tone(4).Minor()
var Am Chord = A.Tone(4).Minor()
var ASm Chord = AS.Tone(4).Minor()
var Bm Chord = B.Tone(4).Minor()

var MINORS []Chord = []Chord{Cm, CSm, Dm, DSm, Em, Fm, FSm, Gm, GSm, Am, ASm, Bm}

// Augmented Triads
// Diminished Triads
// Major 7ths
// Minor 7ths
// Augmented 7ths
// Diminished 7ths
