package music

import (
	"testing"
)

// Major Triads
var C_MAJOR Chord = C.Tone(4).Major()
var CS_MAJOR Chord = CS.Tone(4).Major()
var D_MAJOR Chord = D.Tone(4).Major()
var DS_MAJOR Chord = DS.Tone(4).Major()
var E_MAJOR Chord = E.Tone(4).Major()
var F_MAJOR Chord = F.Tone(4).Major()
var FS_MAJOR Chord = FS.Tone(4).Major()
var G_MAJOR Chord = G.Tone(4).Major()
var GS_MAJOR Chord = GS.Tone(4).Major()
var A_MAJOR Chord = A.Tone(4).Major()
var AS_MAJOR Chord = AS.Tone(4).Major()
var B_MAJOR Chord = B.Tone(4).Major()

var MAJORS []Chord = []Chord{C_MAJOR, CS_MAJOR, D_MAJOR, DS_MAJOR, E_MAJOR, F_MAJOR, FS_MAJOR, G_MAJOR, GS_MAJOR, A_MAJOR, AS_MAJOR, B_MAJOR}

// First Inversion Major Triads
// Second Inversion Major Triads
// Minor Triads
var C_MINOR Chord = C.Tone(4).Minor()
var CS_MINOR Chord = CS.Tone(4).Minor()
var D_MINOR Chord = D.Tone(4).Minor()
var DS_MINOR Chord = DS.Tone(4).Minor()
var E_MINOR Chord = E.Tone(4).Minor()
var F_MINOR Chord = F.Tone(4).Minor()
var FS_MINOR Chord = FS.Tone(4).Minor()
var G_MINOR Chord = G.Tone(4).Minor()
var GS_MINOR Chord = GS.Tone(4).Minor()
var A_MINOR Chord = A.Tone(4).Minor()
var AS_MINOR Chord = AS.Tone(4).Minor()
var B_MINOR Chord = B.Tone(4).Minor()

var MINORS []Chord = []Chord{C_MINOR, CS_MINOR, D_MINOR, DS_MINOR, E_MINOR, F_MINOR, FS_MINOR, G_MINOR, GS_MINOR, A_MINOR, AS_MINOR, B_MINOR}

// Augmented Triads
// Diminished Triads
// Major 7ths
// Minor 7ths
// Augmented 7ths
// Diminished 7ths

var TriadsTestCases []struct {
	Chord
	notes []Tone
} = []struct {
	Chord
	notes []Tone
}{
	{C_MAJOR, []Tone{C.Tone(4), E.Tone(4), G.Tone(4)}},
	{C_MINOR, []Tone{C.Tone(4), DS.Tone(4), G.Tone(4)}},
	{CS_MAJOR, []Tone{CS.Tone(4), F.Tone(4), GS.Tone(4)}},
	{CS_MINOR, []Tone{CS.Tone(4), E.Tone(4), GS.Tone(4)}},
	{D_MAJOR, []Tone{D.Tone(4), FS.Tone(4), A.Tone(4)}},
	{D_MINOR, []Tone{D.Tone(4), F.Tone(4), A.Tone(4)}},
	{DS_MAJOR, []Tone{DS.Tone(4), G.Tone(4), AS.Tone(4)}},
	{DS_MINOR, []Tone{DS.Tone(4), FS.Tone(4), AS.Tone(4)}},
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
	{C_MAJOR, 1, []Tone{E.Tone(4), G.Tone(4), C.Tone(5)}},
	{C_MINOR, 1, []Tone{DS.Tone(4), G.Tone(4), C.Tone(5)}},
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
	{C_MAJOR, []uint8{1, 2, 3, 1}, 4, []Tone{C.Tone(4), E.Tone(4), G.Tone(4), C.Tone(4)}},
}

func TestArpeggiate(t *testing.T) {
	for _, tc := range ArpeggiateTestCases {
		if val := tc.Chord.Arpeggiate(tc.pattern, tc.count); !SameTones(val, tc.tones) {
			t.Errorf("Wrong tones for %s.Arpeggiate(%v, %v), expected: %v got %v", tc.Chord.Name(), tc.pattern, tc.count, ListTones(tc.tones), ListTones(val))
		}
	}
}

func SameTones(s1 []Tone, s2 []Tone) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
