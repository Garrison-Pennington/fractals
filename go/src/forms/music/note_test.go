package music

import "testing"

var AsStringTestCases []struct {
	Note
	string
} = []struct {
	Note
	string
}{
	{C.Note(4), "C4"},
	{CS.Note(4), "C#4"},
	{D.Note(4), "D4"},
	{DS.Note(4), "D#4"},
	{E.Note(4), "E4"},
	{F.Note(4), "F4"},
	{FS.Note(4), "F#4"},
	{G.Note(4), "G4"},
	{GS.Note(4), "G#4"},
	{A.Note(4), "A4"},
	{AS.Note(4), "A#4"},
	{B.Note(4), "B4"},
}

func TestAsString(t *testing.T) {
	for _, tc := range AsStringTestCases {
		if val := tc.Note.AsString(); val != tc.string {
			t.Errorf("Wrong string representation, expected: %s, got: %s", tc.string, val)
		}
	}
}

var OctaveTestCases []struct {
	a Note
	b Note
} = []struct {
	a Note
	b Note
}{
	{A.Note(4), A.Note(5)},
	{AS.Note(4), AS.Note(5)},
	{B.Note(4), B.Note(5)},
	{C.Note(4), C.Note(5)},
	{CS.Note(4), CS.Note(5)},
	{D.Note(4), D.Note(5)},
	{DS.Note(4), DS.Note(5)},
	{E.Note(4), E.Note(5)},
	{F.Note(4), F.Note(5)},
	{FS.Note(4), FS.Note(5)},
	{G.Note(4), G.Note(5)},
	{GS.Note(4), GS.Note(5)},
}

func TestOctaveChange(t *testing.T) {
	for _, tc := range OctaveTestCases {
		if val := tc.a.OctaveUp(1); val != tc.b {
			t.Errorf("Wrong value on OctaveUp, expected: %s, got %s", tc.b.AsString(), val.AsString())
		}
		if val := tc.b.OctaveDown(1); val != tc.a {
			t.Errorf("Wrong value on OctaveDown, expected: %s, got %s", tc.a.AsString(), val.AsString())
		}
	}
}

var PitchChangeTestCases []struct {
	a     Note
	steps uint8
	b     Note
} = []struct {
	a     Note
	steps uint8
	b     Note
}{
	{A.Note(4), 1, AS.Note(4)},   // Minor Second
	{AS.Note(4), 2, C.Note(5)},   // Major Second
	{B.Note(4), 3, D.Note(5)},    // Minor Third
	{C.Note(5), 4, E.Note(5)},    // Major Third
	{CS.Note(5), 5, FS.Note(5)},  // Perfect Fourth
	{D.Note(5), 6, GS.Note(5)},   // Diminished Fifth
	{DS.Note(5), 7, AS.Note(5)},  // Perfect Fifth
	{E.Note(5), 8, C.Note(6)},    // Minor Sixth
	{F.Note(5), 9, D.Note(6)},    // Major Sixth
	{FS.Note(5), 10, E.Note(6)},  // Minor Seventh
	{G.Note(5), 11, FS.Note(6)},  // Major Seventh
	{GS.Note(5), 12, GS.Note(6)}, // Perfect Octave
}

func TestPitchChange(t *testing.T) {
	for _, tc := range PitchChangeTestCases {
		if val := tc.a.PitchedUp(tc.steps); val != tc.b {
			t.Errorf("Wrong value on %s.PitchedUp(%v), expected: %s, got %s", tc.a.AsString(), tc.steps, tc.b.AsString(), val.AsString())
		}
		if val := tc.b.PitchedDown(tc.steps); val != tc.a {
			t.Errorf("Wrong value on %s.PitchedDown(%v), expected: %s, got %s", tc.b.AsString(), tc.steps, tc.a.AsString(), val.AsString())
		}
	}
}
