package examples

import (
	c "fractals/forms/complex"
	"testing"
)

var MANDELBROT Mandelbrot = Mandelbrot{}

var NULL c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(0, 0),
	Computed: complex(0, 0),
	Steps:    0,
}

var NULL_STEPPED c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(0, 0),
	Computed: complex(0, 0),
	Steps:    1,
}

var CV1 c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(1, 1),
	Computed: complex(1, 1),
	Steps:    0,
}

var CV1_STEPPED c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(1, 1),
	Computed: complex(1, 3),
	Steps:    1,
}

var CV2 c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(2, 2),
	Computed: complex(2, 2),
	Steps:    0,
}

var CV2_STEPPED c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(2, 2),
	Computed: complex(2, 10),
	Steps:    1,
}

var STEP_CASES map[c.ComplexFractalValue]c.ComplexFractalValue = map[c.ComplexFractalValue]c.ComplexFractalValue{
	NULL: NULL_STEPPED,
	CV1:  CV1_STEPPED,
	CV2:  CV2_STEPPED,
}

func TestStep(t *testing.T) {
	for check, expect := range STEP_CASES {
		actual := MANDELBROT.Step(check)
		if actual != expect {
			t.Errorf(
				`FIELD: (INPUT, EXPECTED, ACTUAL)
        Initial: (%v, %v, %v)
        Computed: (%v, %v, %v)
        Steps: (%v, %v, %v)`,
				check.Initial, expect.Initial, actual.Initial,
				check.Computed, expect.Computed, actual.Computed,
				check.Steps, expect.Steps, actual.Steps,
			)
		}
	}
}

var NULL_100_STEPS c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(0, 0),
	Computed: complex(0, 0),
	Steps:    100,
}

var CV1_10_STEPS c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(1, 1),
	Computed: complex(1, 3),
	Steps:    1,
}

var CV3 c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(.25, .25),
	Computed: complex(.25, .25),
	Steps:    0,
}

var CV3_3_STEPS c.ComplexFractalValue = c.ComplexFractalValue{
	Initial:  complex(.25, .25),
	Computed: complex(float64(361)/float64(4096), float64(205)/float64(512)),
	Steps:    3,
}

var N_STEP_CASES map[c.ComplexFractalValue]c.ComplexFractalValue = map[c.ComplexFractalValue]c.ComplexFractalValue{
	NULL: NULL_100_STEPS,
	CV1:  CV1_10_STEPS,
	CV3:  CV3_3_STEPS,
}

func TestNSteps(t *testing.T) {
	for check, expect := range N_STEP_CASES {
		actual := c.NSteps(MANDELBROT, check.Initial, expect.Steps)
		if actual != expect {
			t.Errorf(
				`FIELD: (INPUT, EXPECTED, ACTUAL)
        Initial: (%v, %v, %v)
        Computed: (%v, %v, %v)
        Steps: (%v, %v, %v)`,
				check.Initial, expect.Initial, actual.Initial,
				check.Computed, expect.Computed, actual.Computed,
				check.Steps, expect.Steps, actual.Steps,
			)
		}
	}
}
