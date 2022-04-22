package examples

import (
	cmp "fractals/forms/complex"
	"image/color"
	"math"
	"math/cmplx"
)

type DivMod struct {
	Divisor float64
}

func (dm DivMod) Step(val cmp.ComplexFractalValue) cmp.ComplexFractalValue {
	i := math.Mod(imag(val.Computed), dm.Divisor)
	r := real(val.Computed) / dm.Divisor
	c := complex(r, i)
	return val.Step(cmplx.Pow(c, 2) + val.Initial)
}
func (dm DivMod) Exit(val cmp.ComplexFractalValue) bool {
	return cmplx.Abs(val.Computed) >= dm.Divisor
}
func (dm DivMod) Init(c complex128) cmp.ComplexFractalValue {
	return cmp.ComplexFractalValue{
		Initial:  c,
		Computed: c,
		Steps:    0,
	}
}
func (dm DivMod) Color(val cmp.ComplexFractalValue, maxSteps uint32) color.Color {
	score := float64(val.Steps) / float64(maxSteps)
	r := uint8(score * 255)
	g := uint8(50 + score*205)
	b := uint8(255)
	a := uint8(255)
	return color.RGBA{r, g, b, a}
}
func (dm DivMod) Name() string {
	return "DivMod"
}
