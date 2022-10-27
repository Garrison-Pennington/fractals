package music

import (
	"math"
	"math/rand"

	"github.com/rs/zerolog/log"
)

var CHROMATIC_SCALE [12]string = [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
var MAJOR_SCALE_INTERVALS []uint8 = []uint8{2, 2, 1, 2, 2, 2, 1}
var MINOR_SCALE_INTERVALS []uint8 = []uint8{2, 1, 2, 2, 1, 2, 2}

func TransposeOctave(notes []uint8, octaves int8) []uint8 {
	for i := range notes {
		notes[i] = uint8(int8(notes[i]) + 12*octaves)
	}
	return notes
}

type Scale struct {
	Notes [8]string
	Base  uint8 // Numeric component of 'C4', 'C5', etc.
}

func GetScale(base uint8, intervals []uint8) Scale {
	notes := [8]string{}
	var current uint8 = base
	for i, v := range intervals {
		notes[i] = CHROMATIC_SCALE[current%12]
		current += v
	}
	return Scale{notes, base}
}

func (s Scale) translate(nums []uint8) (notes []string) {
	numNotes := uint8(len(nums))
	notes = make([]string, numNotes)
	for i, v := range nums {
		idx := v % 8
		notes[i] = s.Notes[idx]
	}
	return
}

func BridgeSeries(expansions uint8, initiator []uint8) (series []uint8) {
	// Initialize slice with proper size
	numVals := int(math.Pow(float64(len(initiator)), float64(expansions)))
	numInit := len(initiator)
	nums := make([]uint8, numVals)
	// Populate front with initiator
	for exp := uint8(0); exp < expansions; exp++ {
		chunkSize := int(math.Pow(float64(numInit), float64(exp)))
		for idx := 0; idx < numVals; {
			for _, v := range initiator {
				for i := 0; i < chunkSize; i++ {
					nums[idx] += v
					idx++
				}
			}
		}
	}
	return nums
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

func RandomMeasures(n uint8, measures [][]uint16) (ret [][]uint16) {
	right := len(measures)
	ret = make([][]uint16, 0)
	for n > 0 {
		ret = append(ret, measures[rand.Intn(right)])
		n--
	}
	return
}

func ConcatMeasures(measures [][]uint16) (ret []uint16) {
	ret = make([]uint16, 0)
	for _, m := range measures {
		ret = append(ret, m...)
	}
	return
}

// func main() {
// 	Fsm := GetScale(6, MINOR_SCALE_INTERVALS)
// 	notes := Fsm.BridgeSeries(2, []uint8{1, 7, 3, 8})
// 	fmt.Printf("%v: %v", notes, len(notes))
// }
