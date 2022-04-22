package complex

import "fmt"

type Point interface {
	X() float64
	Y() float64
	FromFloats(float64, float64) Point
	FromComplex(complex128) Point
}

func Repr(p Point) string {
	return fmt.Sprintf("(%v4, %v4)", p.X(), p.Y())
}

func PtAdd(p1 Point, p2 Point) Point {
	return p1.FromFloats(p1.X()+p2.X(), p1.Y()+p2.Y())
}

func PtSub(p1 Point, p2 Point) Point {
	return p1.FromFloats(p1.X()-p2.X(), p1.Y()-p2.Y())
}

func PtScale(p1 Point, s float64) Point {
	return p1.FromFloats(p1.X()*s, p1.Y()*s)
}

func PtProd(p1 Point, p2 Point) Point {
	return p1.FromFloats(p1.X()*p2.X(), p1.Y()*p2.Y())
}

func PtDiv(p1 Point, p2 Point) Point {
	return p1.FromFloats(p1.X()/p2.X(), p1.Y()/p2.Y())
}

func AsComplex(pt Point) complex128 {
	return complex(pt.X(), pt.Y())
}

type ComplexPoint struct {
	pt complex128
}

func (c ComplexPoint) X() float64 {
	return real(c.pt)
}

func (c ComplexPoint) Y() float64 {
	return imag(c.pt)
}

func (c ComplexPoint) FromFloats(x float64, y float64) Point {
	return ComplexPoint{complex(x, y)}
}

func (c ComplexPoint) FromComplex(val complex128) Point {
	return ComplexPoint{val}
}

func MakeComplexPoint(c complex128) ComplexPoint {
	return ComplexPoint{c}
}

type Pixel struct {
	x float64
	y float64
}

func (p Pixel) X() float64 {
	return p.x
}

func (p Pixel) Y() float64 {
	return p.y
}

func (p Pixel) PX() uint32 {
	return uint32(p.x)
}

func (p Pixel) PY() uint32 {
	return uint32(p.y)
}

func (p Pixel) FromFloats(x float64, y float64) Point {
	return Pixel{x, y}
}

func (p Pixel) FromComplex(c complex128) Point {
	return Pixel{real(c), imag(c)}
}

func MakePixel(x int, y int) Pixel {
	return Pixel{float64(x), float64(y)}
}

type Transform interface {
	Shape() Point
	Translation() Point
	ConvertPoint(Point) Point
}

func TranslateOut(t Transform, p Point) Point {
	return PtSub(p, t.Translation())
}

func TranslateInto(t Transform, p Point) Point {
	tp := t.Translation()
	p = PtAdd(p, tp)
	return p
}

// Normalize a Vector(Point) to it's own space
// Vector starts at origin and does not exceed length 1 in either dimension
func Normalize(t Transform, p Point) Point {
	p = TranslateOut(t, p)
	return PtDiv(p, t.Shape())
}

// Project a normalized Vector(Point) onto another space
func Project(t Transform, p Point) Point {
	p = PtProd(p, t.Shape())
	p = TranslateInto(t, p)
	return t.ConvertPoint(p)
}

// Transform a point from T1's space to T2's
func TransformPoint(t1 Transform, t2 Transform, p Point) Point {
	p = Normalize(t1, p)

	return Project(t2, p)
}

type Bounds struct {
	Lower Point
	Upper Point
}

func (b Bounds) Repr() string {
	return fmt.Sprintf("L:(%v4,%v4)_R:(%v4,%v4)", b.Lower.X(), b.Lower.Y(), b.Upper.X(), b.Upper.Y())
}

// Return largest step sizes on real and complex axes that render bounds within Resolution
func (b Bounds) Delta(res Resolution) (float64, float64) {
	shape := b.Shape()
	return shape.X() / float64(res.Width), shape.Y() / float64(res.Height)
}

func (b Bounds) Shape() Point {
	return PtSub(b.Upper, b.Lower)
}

func (b Bounds) Translation() Point {
	return b.Lower
}

func (b Bounds) ConvertPoint(p Point) Point {
	return MakeComplexPoint(complex(p.X(), p.Y()))
}

func (b Bounds) Zoom(focalPoint ComplexPoint, scale float64) Bounds {
	newShape := PtDiv(b.Shape(), MakeComplexPoint(complex(scale, scale)))
	centerOffset := PtDiv(newShape, MakeComplexPoint(complex(2, 2)))
	newLower := PtSub(focalPoint, centerOffset)
	newUpper := PtAdd(focalPoint, centerOffset)
	return Bounds{newLower, newUpper}
}

type Resolution struct {
	Height uint32
	Width  uint32
}

var UHD Resolution = Resolution{
	Height: 1080, Width: 1920,
}

func (r Resolution) Repr() string {
	return fmt.Sprintf("%vx%v", r.Width, r.Height)
}

func (r Resolution) Shape() Point {
	return MakePixel(int(r.Width), int(r.Height))
}

func (r Resolution) Translation() Point {
	return Pixel{0, 0}
}

func (r Resolution) ConvertPoint(p Point) Point {
	return MakePixel(int(p.X()), int(p.Y()))
}

// Generator yielding pixels in column major order
func (r Resolution) GenPixels() chan Pixel {
	ch := make(chan Pixel)
	go func() {
		defer close(ch)
		for x := 0; x < int(r.Width); x++ {
			for y := 0; y < int(r.Height); y++ {
				_x := x
				_y := y
				ch <- MakePixel(_x, _y)
			}
		}
	}()
	return ch
}

// Initialize a 2D array with dimensions of Resolution
func (r Resolution) Array() [][]float64 {
	arr := make([][]float64, r.Height)
	for i := range arr {
		arr[i] = make([]float64, r.Width)
	}
	return arr
}

// Initialize a 2D array with dimensions of Resolution
func (r Resolution) ComplexArray() [][]ComplexFractalValue {
	arr := make([][]ComplexFractalValue, r.Height)
	for i := range arr {
		arr[i] = make([]ComplexFractalValue, r.Width)
	}
	return arr
}

func GenFractalArray(arr [][]ComplexFractalValue) chan ComplexFractalValue {
	ch := make(chan ComplexFractalValue)
	go func() {
		defer close(ch)
		for y := range arr {
			for x := range arr[y] {
				ch <- arr[y][x]
			}
		}
	}()
	return ch
}
