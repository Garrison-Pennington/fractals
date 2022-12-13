package music

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
