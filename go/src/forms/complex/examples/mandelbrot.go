package mandelbrot

import (
	cmp "fractals/forms/complex"
	"image/color"
	"math/cmplx"
)

type Mandelbrot struct {
}

func (m Mandelbrot) Name() string {
	return "Mandelbrot"
}

func (m Mandelbrot) Init(c complex128) cmp.ComplexFractalValue {
	return cmp.ComplexFractalValue{Initial: c, Computed: c, Steps: 0}
}

func (m Mandelbrot) Step(c cmp.ComplexFractalValue) cmp.ComplexFractalValue {
	return c.Step(c.Computed*c.Computed + c.Initial)
}

func (m Mandelbrot) Exit(c cmp.ComplexFractalValue) bool {
	return cmplx.Abs(c.Computed) >= 2
}

func (m Mandelbrot) Color(c cmp.ComplexFractalValue, maxSteps uint32) color.Color {
	score := float64(c.Steps) / float64(maxSteps)
	r := uint8(score * 255)
	g := uint8(50 + score*205)
	b := uint8(255)
	a := uint8(255)
	return color.RGBA{r, g, b, a}
}
