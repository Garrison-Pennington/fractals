package music

import (
	"math"
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

// func main() {
// 	Fsm := GetScale(6, MINOR_SCALE_INTERVALS)
// 	notes := Fsm.BridgeSeries(2, []uint8{1, 7, 3, 8})
// 	fmt.Printf("%v: %v", notes, len(notes))
// }
