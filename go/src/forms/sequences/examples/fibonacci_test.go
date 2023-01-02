package fibonacci

import (
	seq "fractals/forms/sequences"
	"testing"
)

func TestNextN(t *testing.T) {
	expect := [10]uint{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	vals := seq.NextUints(&BASE_FIB, 10)
	for i, val := range vals {
		if val != expect[i] {
			t.Errorf("Fibonacci value #%v incorrect, expected: %v, got: %v", i+1, expect[i], val)
		}
	}
}
