package music_theory

import "testing"

// Base Measures
var M_3_4 Measure = NewMeasure(SIG_3_4)

// Tone Sequences
var C_ARP_ASC []Tone = CM.Arpeggiate([]uint8{1, 2, 3}, 3)

// Note Sequences
var C_ARP_ASC_WHOLES []Note = AsWholes(C_ARP_ASC...)
var C_ARP_ASC_QUARTERS []Note = AsQuarters(C_ARP_ASC...)

var AddNotesTestCases []struct {
	Measure
	notes []Note
} = []struct {
	Measure
	notes []Note
}{
	{NewMeasure(SIG_4_4), C_ARP_ASC_WHOLES},
	{NewMeasure(SIG_3_4), C_ARP_ASC_QUARTERS},
}

func TestAddNotes(t *testing.T) {
	for tn, tc := range AddNotesTestCases {
		tc.Measure.AddNotes(0, tc.notes...)
		last := 0
		for _, notes := range tc.Measure.Notes {
			for _, check := range notes {
				if check != tc.notes[last] {
					t.Errorf("Wrong notes in Measure.AddNotes case #%v expected: %v, got %v", tn, ListNotes(tc.notes), ListNotes(notes))
					break
				}
				last++
			}
		}
	}
}

var MeasureAsStringTestCases []struct {
	Measure
	notes  []Note
	expect string
} = []struct {
	Measure
	notes  []Note
	expect string
}{
	{NewMeasure(SIG_3_4), C_ARP_ASC_QUARTERS, "0. C4Q \n8. E4Q \n16. G4Q \n"},
}

func TestMeasureAsString(t *testing.T) {
	for tn, tc := range MeasureAsStringTestCases {
		tc.Measure.AddNoteSequence(0, tc.notes...)
		if val := tc.Measure.AsString(); val != tc.expect {
			t.Errorf("Wrong string rep for Measure.AsString() case #%v expected: %s, got %s", tn, tc.expect, val)
		}
	}
}
