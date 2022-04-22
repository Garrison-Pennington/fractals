package main

import (
	cmp "fractals/forms/complex"
	mand "fractals/forms/complex/examples"
)

var MANDELBROT mand.Mandelbrot = mand.Mandelbrot{}

var UHD cmp.Resolution = cmp.Resolution{
	Width: 1920, Height: 1080,
}

var MANDELBROT_WHOLE_BOUNDS cmp.Bounds = cmp.Bounds{
	Lower: cmp.MakeComplexPoint(complex(-2, -2)),
	Upper: cmp.MakeComplexPoint(complex(2, 2)),
}
var MANDELBROT_REDUCED_BOUNDS cmp.Bounds = cmp.Bounds{
	Lower: cmp.MakeComplexPoint(complex(-2, -1.5)),
	Upper: cmp.MakeComplexPoint(complex(1.5, 1.5)),
}

var DEFAULT_STEPS uint32 = 100
var DEEP_STEPS uint32 = 1000

var UHD_MANDELBROT_FULL = cmp.ComplexRender{
	Resolution: UHD,
	Bounds:     MANDELBROT_WHOLE_BOUNDS,
	Fractal:    MANDELBROT,
	Steps:      DEFAULT_STEPS,
}

var WINDOW_10X10 cmp.Window = cmp.Window{
	Height: 10,
	Width:  10,
	Array:  nil,
}

var STEP_PROP_DEFAULT func(cmp.ComplexFractalValue) float64 = cmp.StepScorer(DEFAULT_STEPS)
var REACHED_STEPS_DEFAULT func(cmp.ComplexFractalValue) bool = cmp.BinStepScorer(DEFAULT_STEPS)

func main() {
	UHD_MANDELBROT_FULL.MultiPhaseRender(5, WINDOW_10X10, 5).Save()

}
