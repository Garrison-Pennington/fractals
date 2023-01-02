package music_theory

import (
	"testing"
)

var TriadsTestCases []struct {
	Chord
	notes []Tone
} = []struct {
	Chord
	notes []Tone
}{
	{CM, []Tone{C.Tone(4), E.Tone(4), G.Tone(4)}},
	{Cm, []Tone{C.Tone(4), DS.Tone(4), G.Tone(4)}},
	{CSM, []Tone{CS.Tone(4), F.Tone(4), GS.Tone(4)}},
	{CSm, []Tone{CS.Tone(4), E.Tone(4), GS.Tone(4)}},
	{DM, []Tone{D.Tone(4), FS.Tone(4), A.Tone(4)}},
	{Dm, []Tone{D.Tone(4), F.Tone(4), A.Tone(4)}},
	{DSM, []Tone{DS.Tone(4), G.Tone(4), AS.Tone(4)}},
	{DSm, []Tone{DS.Tone(4), FS.Tone(4), AS.Tone(4)}},
}

func TestTriads(t *testing.T) {
	for _, tc := range TriadsTestCases {
		if !SameTones(tc.Chord.Tones, tc.notes) {
			t.Errorf("Wrong notes for %s, expected: %s got: %s", tc.Chord.Name(), ListTones(tc.notes), ListTones(tc.Chord.Tones))
		}
	}
}

var InvertTestCases []struct {
	Chord
	inversions uint8
	notes      []Tone
} = []struct {
	Chord
	inversions uint8
	notes      []Tone
}{
	{CM, 1, []Tone{E.Tone(4), G.Tone(4), C.Tone(5)}},
	{Cm, 1, []Tone{DS.Tone(4), G.Tone(4), C.Tone(5)}},
}

func TestInvert(t *testing.T) {
	for _, tc := range InvertTestCases {
		tc.Chord.Invert()
		if !SameTones(tc.Chord.Tones, tc.notes) {
			t.Errorf("%s.Invert() got wrong notes expected: %s got: %s", tc.Chord.Name(), ListTones(tc.notes), ListTones(tc.Chord.Tones))
		}
	}
}

var ArpeggiateTestCases []struct {
	Chord
	pattern []uint8
	count   uint8
	tones   []Tone
} = []struct {
	Chord
	pattern []uint8
	count   uint8
	tones   []Tone
}{
	{CM, []uint8{1, 2, 3, 1}, 4, []Tone{C.Tone(4), E.Tone(4), G.Tone(4), C.Tone(4)}},
	{CM, []uint8{1, 2, 3, 1}, 8, []Tone{C.Tone(4), E.Tone(4), G.Tone(4), C.Tone(4), C.Tone(4), E.Tone(4), G.Tone(4), C.Tone(4)}},
}

func TestArpeggiate(t *testing.T) {
	for _, tc := range ArpeggiateTestCases {
		if val := tc.Chord.Arpeggiate(tc.pattern, tc.count); !SameTones(val, tc.tones) {
			t.Errorf("Wrong tones for %s.Arpeggiate(%v, %v), expected: %v got %v", tc.Chord.Name(), tc.pattern, tc.count, ListTones(tc.tones), ListTones(val))
		}
	}
}
