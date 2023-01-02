package sequences

type SequenceParameters struct {
	Floats  map[string]float64
	Ints    map[string]int
	Strings map[string]string
	Keys    map[string][]string
}

type FloatSequence interface {
	Gen() chan float64
	Next() float64
}

type Uint8Sequence interface {
	Gen() chan uint8
	Next() uint8
}

type IntSequence interface {
	Gen() chan int
	Next() int
}

type UintSequence interface {
	Gen() chan uint
	Next() uint
}

func NextFloats(s FloatSequence, n int) (nums []float64) {
	for i := 0; i < n; i++ {
		nums = append(nums, s.Next())
	}
	return
}

func NextInts(s IntSequence, n int) (nums []int) {
	for i := 0; i < n; i++ {
		nums = append(nums, s.Next())
	}
	return
}

func NextUints(s UintSequence, n int) (nums []uint) {
	for i := 0; i < n; i++ {
		nums = append(nums, s.Next())
	}
	return
}
