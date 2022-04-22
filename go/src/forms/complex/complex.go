package complex

import "image/color"

type ComplexFractal interface {
	Init(complex128) ComplexFractalValue           // Prepare initial value of equation using provided complex number
	Step(ComplexFractalValue) ComplexFractalValue  // Perform one iteration of the fractal equation
	Exit(ComplexFractalValue) bool                 // Should the iteration stop and return a final value?
	Color(ComplexFractalValue, uint32) color.Color // Return a color corresponding to the final value
	Name() string
}

type ComplexFractalValue struct {
	Initial  complex128 // Variable component of fractal equation
	Computed complex128 // Resulting value after Steps iterations
	Steps    uint32     // Number of iterations to produce Computed
}

func (c ComplexFractalValue) Step(z complex128) ComplexFractalValue {
	return ComplexFractalValue{c.Initial, z, c.Steps + 1}
}

// Iterate a ComplexFractal n times, using the output of the last iteration as the arg to the next
func NSteps(fractal ComplexFractal, c complex128, n uint32) ComplexFractalValue {
	value := fractal.Init(c)
	for n > 0 {
		value = fractal.Step(value)
		if fractal.Exit(value) {
			break
		}
		n--
	}
	return value
}
