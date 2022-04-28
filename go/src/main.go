package main

import (
	"flag"
	cmp "fractals/forms/complex"
	mand "fractals/forms/complex/examples"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	resPtr := flag.String("resolution", "UHD", "Render output resolution")
	fractalPtr := flag.String("fractal", "mandelbrot", "The fractal to render")
	stepsPtr := flag.Uint("steps", 100, "Number of iterations to run on each point")
	phasesPtr := flag.Uint("phases", 0, "Number of zoom phases to perform in MultiPhaseRender, 0 is regular render")
	windowPtr := flag.Int("window-size", 10, "Side length of the MultiPhaseWindow")
	stridePtr := flag.Int("stride", 5, "Stride of the MultiPhaseWindow")
	flag.Parse()
	args := os.Args[1:]
	switch action := args[0]; action {
	case "render":
		render := cmp.ComplexRender{
			Resolution: parseResolution(*resPtr),
			Bounds:     parseBounds(*fractalPtr),
			Fractal:    parseFractal(*fractalPtr),
			Steps:      uint32(*stepsPtr),
		}
		if *phasesPtr > 0 {
			window := cmp.Window{
				Height: *windowPtr,
				Width:  *windowPtr,
				Array:  nil,
			}
			render.MultiPhaseRender(uint8(*phasesPtr), window, *stridePtr).Save()
		} else {
			render.Save()
		}
		break
	case "music":
		break
	}

	//UHD_MANDELBROT_FULL.MultiPhaseRender(5, WINDOW_10X10, 5).Save()

}

func parseResolution(resStr string) cmp.Resolution {
	switch resStr {
	case "UHD":
		return UHD
	}
	log.Fatal().Msgf("Invalid render resolution given: %v", resStr)
	return UHD // TODO: Make nil value
}

func parseBounds(fractalStr string) cmp.Bounds {
	switch fractalStr {
	case "mandelbrot":
		return MANDELBROT_WHOLE_BOUNDS
	case "divmod":
		return cmp.Bounds{
			Lower: cmp.MakeComplexPoint(complex(-5, -5)),
			Upper: cmp.MakeComplexPoint(complex(5, 5)),
		}
	}
	log.Fatal().Msgf("Invalid render bounds given: %v", fractalStr)
	return MANDELBROT_WHOLE_BOUNDS // TODO: Make nil value
}

func parseFractal(fractalStr string) cmp.ComplexFractal {
	switch fractalStr {
	case "mandelbrot":
		return MANDELBROT
	case "divmod":
		break // TODO: Merge new-fractal branch and add DivMod here
	}
	log.Fatal().Msgf("Invalid render resolution given: %v", fractalStr)
	return MANDELBROT // TODO: Make nil value
}
