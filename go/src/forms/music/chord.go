package music

type Chord struct {
	Root    Tone
	Notes   []Note
	Quality string
}

// CONSTRUCTORS:
func MinorTriad(root Note) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root)}, "m"}
}

func MajorTriad(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root)}, ""}
}

func AugmentedTriad(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), AUGMENTED_5TH.NextNote(root)}, "+"}
}

func DiminishedTriad(root Note) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), DIMINISHED_5TH.NextNote(root)}, "o"}
}

func DiminishedSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), DIMINISHED_5TH.NextNote(root), DIMINISHED_7TH.NextNote(root)}, "o7"}
}

func HalfDiminishedSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), DIMINISHED_5TH.NextNote(root), MINOR_7TH.NextNote(root)}, root.Tone.Value}
}

func MinorSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root), MINOR_7TH.NextNote(root)}, "m7"}
}

func MinorMajorSeventh(root Note, inversions uint8) Chord {
	return Chord{root.Tone, []Note{root, MINOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root), MAJOR_7TH.NextNote(root)}, root.Tone.Value}
}

func DominantSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root), MINOR_7TH.NextNote(root)}, root.Tone.Value + "7"}
}

func MajorSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), PERFECT_5TH.NextNote(root), MAJOR_7TH.NextNote(root)}, root.Tone.Value}
}

func AugmentedSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), AUGMENTED_5TH.NextNote(root), MINOR_7TH.NextNote(root)}, root.Tone.Value}
}

func AugmentedMajorSeventh(root Note) Chord {
	return Chord{root.Tone, []Note{root, MAJOR_3RD.NextNote(root), AUGMENTED_5TH.NextNote(root), MAJOR_7TH.NextNote(root)}, root.Tone.Value}
}

// METHODS:
func (c *Chord) Invert() {
	c.Notes = append(c.Notes[1:], c.Notes[0].OctaveUp(1))
}

func (c Chord) Name() string {
	return c.Root.Value + c.Quality
}
