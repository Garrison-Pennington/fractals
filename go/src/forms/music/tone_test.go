package music

import "testing"

var AsStringTestCases []struct {
	Tone
	string
} = []struct {
	Tone
	string
}{
	{C.Tone(4), "C4"},
	{CS.Tone(4), "C#4"},
	{D.Tone(4), "D4"},
	{DS.Tone(4), "D#4"},
	{E.Tone(4), "E4"},
	{F.Tone(4), "F4"},
	{FS.Tone(4), "F#4"},
	{G.Tone(4), "G4"},
	{GS.Tone(4), "G#4"},
	{A.Tone(4), "A4"},
	{AS.Tone(4), "A#4"},
	{B.Tone(4), "B4"},
}

func TestAsString(t *testing.T) {
	for _, tc := range AsStringTestCases {
		if val := tc.Tone.AsString(); val != tc.string {
			t.Errorf("Wrong string representation, expected: %s, got: %s", tc.string, val)
		}
	}
}

var OctaveTestCases []struct {
	a Tone
	b Tone
} = []struct {
	a Tone
	b Tone
}{
	{A.Tone(4), A.Tone(5)},
	{AS.Tone(4), AS.Tone(5)},
	{B.Tone(4), B.Tone(5)},
	{C.Tone(4), C.Tone(5)},
	{CS.Tone(4), CS.Tone(5)},
	{D.Tone(4), D.Tone(5)},
	{DS.Tone(4), DS.Tone(5)},
	{E.Tone(4), E.Tone(5)},
	{F.Tone(4), F.Tone(5)},
	{FS.Tone(4), FS.Tone(5)},
	{G.Tone(4), G.Tone(5)},
	{GS.Tone(4), GS.Tone(5)},
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
	a     Tone
	steps uint8
	b     Tone
} = []struct {
	a     Tone
	steps uint8
	b     Tone
}{
	{A.Tone(4), 1, AS.Tone(4)},   // Minor Second
	{AS.Tone(4), 2, C.Tone(5)},   // Major Second
	{B.Tone(4), 3, D.Tone(5)},    // Minor Third
	{C.Tone(5), 4, E.Tone(5)},    // Major Third
	{CS.Tone(5), 5, FS.Tone(5)},  // Perfect Fourth
	{D.Tone(5), 6, GS.Tone(5)},   // Diminished Fifth
	{DS.Tone(5), 7, AS.Tone(5)},  // Perfect Fifth
	{E.Tone(5), 8, C.Tone(6)},    // Minor Sixth
	{F.Tone(5), 9, D.Tone(6)},    // Major Sixth
	{FS.Tone(5), 10, E.Tone(6)},  // Minor Seventh
	{G.Tone(5), 11, FS.Tone(6)},  // Major Seventh
	{GS.Tone(5), 12, GS.Tone(6)}, // Perfect Octave
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
