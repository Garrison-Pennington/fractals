package complex

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

type ComplexRender struct {
	Resolution Resolution
	Bounds     Bounds
	Fractal    ComplexFractal
	Steps      uint32
}

func (r ComplexRender) Render() [][]ComplexFractalValue {
	return Render(r.Fractal, r.Resolution, r.Bounds, r.Steps)
}

func (r ComplexRender) RenderImage() *image.RGBA {
	return RenderImage(r.Fractal, r.Resolution, r.Bounds, r.Steps)
}

func (r ComplexRender) Save() {
	filename := r.Fractal.Name() + "_" + r.Resolution.Repr() + "_" + r.Bounds.Repr() + ".png"
	img := RenderImage(r.Fractal, r.Resolution, r.Bounds, r.Steps)
	SaveImage(filename, img)
}

func (r ComplexRender) Zoom(pt ComplexPoint, scale float64) ComplexRender {
	log.Debug().Msg("Zooming Render")
	return ComplexRender{
		Resolution: r.Resolution,
		Bounds:     r.Bounds.Zoom(pt, scale),
	}
}

func (r ComplexRender) MultiPhaseRender(numPhases uint8, w Window, stride int) ComplexRender {
	last := r
	nextBounds := r.Bounds
	for numPhases > 0 {
		numPhases--
		w = w.IntoArray(last.Render())
		maxEdges := 0
		for view := range w.Slide(stride) {
			edges := view.EdgeCount(r.Steps)
			if edges > maxEdges {
				maxEdges = edges
				nextBounds = view.Bounds()
			}
		}
		last = UHDRender(r.Fractal, nextBounds, r.Steps)
	}
	return last
}

func UHDRender(f ComplexFractal, b Bounds, n uint32) ComplexRender {
	log.Debug().Msgf("Creating UHD Render")
	return ComplexRender{
		Resolution: UHD,
		Bounds:     b,
		Fractal:    f,
		Steps:      n,
	}
}

func Render(fractal ComplexFractal, outputResolution Resolution, bounds Bounds, steps uint32) (arr [][]ComplexFractalValue) {
	log.Debug().Msg("Preparing to Render")
	log.Trace().Msg("Creating 2D array")
	arr = outputResolution.ComplexArray()
	log.Trace().Msg("Creating Pixel Generator")
	ch := outputResolution.GenPixels()
	count := 0
	lastPercent := 0
	totalPixels := outputResolution.Height * outputResolution.Width
	log.Info().Msgf("Beginning Render with %v pixels", totalPixels)
	for pixel := range ch {
		count++
		if percentComplete := (float64(count) / float64(totalPixels)) / .25; int(percentComplete) > lastPercent {
			lastPercent = int(percentComplete)
			log.Info().Msgf("%v%% Complete", lastPercent*25)
		}
		pt := TransformPoint(outputResolution, bounds, pixel)
		arr[pixel.PY()][pixel.PX()] = NSteps(fractal, AsComplex(pt), steps)
	}
	log.Info().Msg("Render Complete")
	return
}

func RenderImage(fractal ComplexFractal, outputResolution Resolution, bounds Bounds, steps uint32) *image.RGBA {
	render := Render(fractal, outputResolution, bounds, steps)
	log.Debug().Msg("Creating Image Background")
	img := Background(outputResolution, color.Black)
	log.Debug().Msg("Drawing Fractal")
	for r, row := range render {
		for c, val := range row {
			img.Set(c, r, fractal.Color(val, steps))
		}
	}
	return img
}

func Background(outputResolution Resolution, c color.Color) *image.RGBA {
	h, w := outputResolution.Height, outputResolution.Width
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	ch := outputResolution.GenPixels()
	for pixel := range ch {
		img.Set(int(pixel.x), int(pixel.y), c)
	}
	return img
}

func SaveImage(filename string, img *image.RGBA) string {
	if !strings.HasSuffix(filename, ".png") {
		filename += ".png"
	}
	log.Info().Msgf("Saving image at %v", filename)
	f, err := os.Create(ImagePath(filename))
	defer f.Close()
	if err != nil {

		log.Fatal().Msg("Error Creating File")
	}
	err = png.Encode(f, img)
	return filename
}

func ImagePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Info().Msg("Error getting current User")
	}
	return path.Join(usr.HomeDir, "src", "fractals", "images", filename)
}
