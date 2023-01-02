package chord_cycle

import (
	seq "fractals/forms/music/sequences"
	mt "fractals/forms/music/theory"
	"testing"
)

var NextChordsTestCases []struct {
	cycle  ChordCycle
	expect []mt.Chord
	count  int
} = []struct {
	cycle  ChordCycle
	expect []mt.Chord
	count  int
}{
	{CM_Em_GM_1, []mt.Chord{mt.CM, mt.Em, mt.GM, mt.CM, mt.Em, mt.GM}, 6},
	{Fm_Cm_CSM_ASm_3, []mt.Chord{mt.Fm, mt.Fm, mt.Fm, mt.Cm, mt.Cm, mt.Cm, mt.CSM, mt.CSM, mt.CSM, mt.ASm, mt.ASm, mt.ASm}, 12},
}

func TestNext(t *testing.T) {
	for _, tc := range NextChordsTestCases {
		if chords := seq.NextChords(&tc.cycle, tc.count); !mt.SameChords(chords, tc.expect) {
			t.Errorf("Wrong chords for %s | Expected: %v | Got: %v", tc.cycle.AsString(), mt.ListChords(tc.expect), mt.ListChords(chords))
		}
	}
}
