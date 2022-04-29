package music

import (
	"github.com/rs/zerolog/log"
	midi "gitlab.com/gomidi/midi/v2"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

func PlayNoteMsgs(channel uint8, note uint8, velocity uint8) (midi.Message, midi.Message) {
	return midi.NoteOn(channel, note, velocity), midi.NoteOff(channel, note)
}

func PlayNote(tr *smf.Track, deltaTicks uint32, channel uint8, note uint8, velocity uint8) {
	on, off := PlayNoteMsgs(channel, note, velocity)
	tr.Add(0, on)
	tr.Add(deltaTicks, off)
}

func Play16th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks16th(), channel, note, velocity)
}

func Play8th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks8th(), channel, note, velocity)
}

func Play4th(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th(), channel, note, velocity)
}

func PlayHalf(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th()*2, channel, note, velocity)
}

func PlayWhole(tr *smf.Track, clock smf.MetricTicks, channel uint8, note uint8, velocity uint8) {
	PlayNote(tr, clock.Ticks4th()*4, channel, note, velocity)
}

type TimeSignature struct {
	Beats uint8
	Value uint8
}

func (ts TimeSignature) MeasureCombos(denom uint16) [][]uint16 {
	total := uint16(ts.Beats) * (denom / uint16(ts.Value))
	cache := make(map[uint16][][]uint16)
	log.Debug().Msgf("Filling %v:%v meter with %vth notes", ts.Beats, ts.Value, denom)
	return ts.FillMeasure(total, denom, cache)
}

func (ts TimeSignature) FillMeasure(needed uint16, denom uint16, cache map[uint16][][]uint16) (combos [][]uint16) {
	log.Debug().Msgf("Filling %v %vth notes", needed, denom)
	if val, ok := cache[needed]; ok {
		log.Debug().Msg("Using cached result")
		return val
	}
	log.Debug().Msg("Calculating combos")
	leaves := make([]uint16, 0)
	complete := make([][]uint16, 0)
	for currDenom := denom; currDenom >= 1; currDenom /= 2 {
		if needed == currDenom {
			complete = append(complete, []uint16{currDenom})
		} else if needed > currDenom {
			leaves = append(leaves, currDenom)
		}
	}
	for _, leaf := range leaves {
		for _, combo := range ts.FillMeasure(needed-leaf, denom, cache) {
			combo = append([]uint16{leaf}, combo...)
			complete = append(complete, combo)
		}
	}
	cache[needed] = complete
	return complete
}
