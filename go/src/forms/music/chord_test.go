package music

import (
	"testing"
)

// Major Triads
var C_MAJOR Chord = C.Note(4).Major()
var CS_MAJOR Chord = CS.Note(4).Major()
var D_MAJOR Chord = D.Note(4).Major()
var DS_MAJOR Chord = DS.Note(4).Major()
var E_MAJOR Chord = E.Note(4).Major()
var F_MAJOR Chord = F.Note(4).Major()
var FS_MAJOR Chord = FS.Note(4).Major()
var G_MAJOR Chord = G.Note(4).Major()
var GS_MAJOR Chord = GS.Note(4).Major()
var A_MAJOR Chord = A.Note(4).Major()
var AS_MAJOR Chord = AS.Note(4).Major()
var B_MAJOR Chord = B.Note(4).Major()

var MAJORS []Chord = []Chord{C_MAJOR, CS_MAJOR, D_MAJOR, DS_MAJOR, E_MAJOR, F_MAJOR, FS_MAJOR, G_MAJOR, GS_MAJOR, A_MAJOR, AS_MAJOR, B_MAJOR}

// First Inversion Major Triads
// Second Inversion Major Triads
// Minor Triads
var C_MINOR Chord = C.Note(4).Minor()
var CS_MINOR Chord = CS.Note(4).Minor()
var D_MINOR Chord = D.Note(4).Minor()
var DS_MINOR Chord = DS.Note(4).Minor()
var E_MINOR Chord = E.Note(4).Minor()
var F_MINOR Chord = F.Note(4).Minor()
var FS_MINOR Chord = FS.Note(4).Minor()
var G_MINOR Chord = G.Note(4).Minor()
var GS_MINOR Chord = GS.Note(4).Minor()
var A_MINOR Chord = A.Note(4).Minor()
var AS_MINOR Chord = AS.Note(4).Minor()
var B_MINOR Chord = B.Note(4).Minor()

var MINORS []Chord = []Chord{C_MINOR, CS_MINOR, D_MINOR, DS_MINOR, E_MINOR, F_MINOR, FS_MINOR, G_MINOR, GS_MINOR, A_MINOR, AS_MINOR, B_MINOR}

// Augmented Triads
// Diminished Triads
// Major 7ths
// Minor 7ths
// Augmented 7ths
// Diminished 7ths

var TriadsTestCases []struct {
	Chord
	notes []Note
} = []struct {
	Chord
	notes []Note
}{
	{C_MAJOR, []Note{C.Note(4), E.Note(4), G.Note(4)}},
	{C_MINOR, []Note{C.Note(4), DS.Note(4), G.Note(4)}},
	{CS_MAJOR, []Note{CS.Note(4), F.Note(4), GS.Note(4)}},
	{CS_MINOR, []Note{CS.Note(4), E.Note(4), GS.Note(4)}},
	{D_MAJOR, []Note{D.Note(4), FS.Note(4), A.Note(4)}},
	{D_MINOR, []Note{D.Note(4), F.Note(4), A.Note(4)}},
	{DS_MAJOR, []Note{DS.Note(4), G.Note(4), AS.Note(4)}},
	{DS_MINOR, []Note{DS.Note(4), FS.Note(4), AS.Note(4)}},
}

func TestTriads(t *testing.T) {
	for _, tc := range TriadsTestCases {
		if !SameNotes(tc.Chord.Notes, tc.notes) {
			t.Errorf("Wrong notes for %s, expected: %s got: %s", tc.Chord.Name(), ListNotes(tc.notes), ListNotes(tc.Chord.Notes))
		}
	}
}

var InvertTestCases []struct {
	Chord
	inversions uint8
	notes      []Note
} = []struct {
	Chord
	inversions uint8
	notes      []Note
}{
	{C_MAJOR, 1, []Note{E.Note(4), G.Note(4), C.Note(5)}},
	{C_MINOR, 1, []Note{DS.Note(4), G.Note(4), C.Note(5)}},
}

func TestInvert(t *testing.T) {
	for _, tc := range InvertTestCases {
		tc.Chord.Invert()
		if !SameNotes(tc.Chord.Notes, tc.notes) {
			t.Errorf("%s.Invert() got wrong notes expected: %s got: %s", tc.Chord.Name(), ListNotes(tc.notes), ListNotes(tc.Chord.Notes))
		}
	}
}

func SameNotes(s1 []Note, s2 []Note) bool {
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
