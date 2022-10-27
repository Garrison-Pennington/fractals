package music

import (
	"math"

	"github.com/rs/zerolog/log"
	reg "github.com/sajari/regression"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

var NOTE_VALUES map[uint16]func(*smf.Track, smf.MetricTicks, uint8, uint8, uint8) = map[uint16]func(*smf.Track, smf.MetricTicks, uint8, uint8, uint8){
	1:  Play16th,
	2:  Play8th,
	4:  Play4th,
	8:  PlayHalf,
	16: PlayWhole,
}

func PinkSlope(noteValues []uint16) float64 {
	bins := make(map[uint16]uint32)
	for _, val := range noteValues {
		bins[val] += 1
	}
	r := new(reg.Regression)
	r.SetObserved("log(frequency)")
	r.SetVar(0, "log(power)")
	dataPts := make(reg.DataPoints, 0)
	for k, v := range bins {
		dataPts = append(dataPts, reg.DataPoint(math.Log(float64(v)), []float64{math.Log(float64(k))}))
	}
	r.Train(dataPts...)
	r.Run()
	log.Info().Msgf("Regression formula:\n%v\n", r.Formula)
	log.Info().Msgf("Regression:\n%s\n", r)
	return 0
}

func MakePink(noteValues []uint16) map[uint16]uint32 {
	bins := make(map[uint16]uint32)
	for _, val := range noteValues {
		bins[val] += 1
	}
	return bins
}
