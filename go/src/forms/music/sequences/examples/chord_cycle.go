package chord_cycle

import (
	mt "fractals/forms/music/theory"
	"strconv"
)

type ChordCycle struct {
	Chords      []mt.Chord
	Repetitions uint8
	rep         uint8
	idx         uint8
}

func MakeChordCycle(chords []mt.Chord, repetitions uint8) ChordCycle {
	return ChordCycle{chords, repetitions, 0, 0}
}

func (cycle *ChordCycle) Next() (chord mt.Chord) {
	cycle.incrementRep()
	chord = cycle.currentChord()
	return
}

func (cycle ChordCycle) currentChord() mt.Chord {
	return cycle.Chords[cycle.idx]
}

func (cycle *ChordCycle) incrementRep() {
	if cycle.rep == cycle.Repetitions {
		cycle.rep = 0
		cycle.incrementIdx()
	}
	cycle.rep++
}

func (cycle *ChordCycle) incrementIdx() {
	if cycle.idx == uint8(len(cycle.Chords)-1) {
		cycle.idx = 0
	} else {
		cycle.idx++
	}
}

func (cycle ChordCycle) AsString() string {
	return mt.ListChords(cycle.Chords) + strconv.FormatInt(int64(cycle.Repetitions), 10)
}

var CM_Em_GM_1 ChordCycle = MakeChordCycle([]mt.Chord{mt.CM, mt.Em, mt.GM}, 1)
var Fm_Cm_CSM_ASm_3 ChordCycle = MakeChordCycle([]mt.Chord{mt.Fm, mt.Cm, mt.CSM, mt.ASm}, 3)
