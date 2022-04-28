package music

import gm "gitlab.com/gomidi/midi/v2/gm"

type Instrument struct {
	Name string
	GM   uint8
}

var PIANO Instrument = Instrument{gm.Instr_AcousticGrandPiano.String(), gm.Instr_AcousticGrandPiano.Value()}
var CELLO Instrument = Instrument{gm.Instr_Cello.String(), gm.Instr_Cello.Value()}
