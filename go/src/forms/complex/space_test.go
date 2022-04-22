package complex

import (
	"math"
	"testing"
)

var ORIGIN_PIXEL Pixel = MakePixel(0, 0)
var TWO_BY_TWO Pixel = MakePixel(2, 2)
var UHD_BOTTOM_RIGHT Pixel = MakePixel(1920, 1080)
var UHD_CENTER Pixel = MakePixel(960, 540)
var MINI_BOTTOM_RIGHT Pixel = MakePixel(100, 100)
var MINI_CENTER Pixel = MakePixel(50, 50)

var ORIGIN_COMPLEX ComplexPoint = MakeComplexPoint(complex(0, 0))
var MANDELBROT_LOWER_BOUND ComplexPoint = MakeComplexPoint(complex(-2, -2))
var MANDELBROT_UPPER_BOUND ComplexPoint = MakeComplexPoint(complex(2, 2))
var MANDELBROT_RANGE ComplexPoint = MakeComplexPoint(complex(4, 4))
var MANDELBROT_BOUND_DIV ComplexPoint = MakeComplexPoint(complex(-1, -1))
var MANDELBROT_ARB_1 ComplexPoint = MakeComplexPoint(complex(-.75, .34))
var MANDELBROT_ARB_2 ComplexPoint = MakeComplexPoint(complex(1.2, -.6))
var MA1_TRANSLATED_OUT ComplexPoint = MakeComplexPoint(complex(1.25, 2.34))

var MINI Resolution = Resolution{Height: 100, Width: 100}

var MANDELBROT_BOUNDS Bounds = Bounds{MANDELBROT_LOWER_BOUND, MANDELBROT_UPPER_BOUND}

// Values are in order: P1 + P2 = Expected
var PT_ADD_CASES = [][3]Point{
	[3]Point{MANDELBROT_LOWER_BOUND, MANDELBROT_UPPER_BOUND, ORIGIN_COMPLEX},
	[3]Point{UHD_CENTER, UHD_CENTER, UHD_BOTTOM_RIGHT},
	[3]Point{MINI_CENTER, MINI_CENTER, MINI_BOTTOM_RIGHT},
}

// Values are in order: P1 - P2 = Expected
var PT_SUB_CASES = [][3]Point{
	[3]Point{MANDELBROT_UPPER_BOUND, MANDELBROT_LOWER_BOUND, MANDELBROT_RANGE},
	[3]Point{UHD_BOTTOM_RIGHT, UHD_CENTER, UHD_CENTER},
	[3]Point{MINI_BOTTOM_RIGHT, MINI_CENTER, MINI_CENTER},
}

// Values are in order: P1 / P2 = Expected
var PT_DIV_CASES = [][3]Point{
	[3]Point{MANDELBROT_UPPER_BOUND, MANDELBROT_LOWER_BOUND, MANDELBROT_BOUND_DIV},
	[3]Point{UHD_BOTTOM_RIGHT, UHD_CENTER, TWO_BY_TWO},
	[3]Point{MINI_BOTTOM_RIGHT, MINI_CENTER, TWO_BY_TWO},
}

func TestPtOps(t *testing.T) {
	arithmeticFuncs := []func(Point, Point) Point{
		PtAdd, PtSub, PtDiv, //PtProd
	}
	arithmeticTestCases := [][][3]Point{
		PT_ADD_CASES, PT_SUB_CASES, PT_DIV_CASES,
	}
	messages := []string{"+", "-", "/"}
	for i, tests := range arithmeticTestCases {
		for _, tc := range tests {
			p1, p2, expect := tc[0], tc[1], tc[2]
			ptFn := arithmeticFuncs[i]
			if actual := ptFn(p1, p2); actual != expect {
				t.Errorf("(%v, %v) %v (%v, %v) = (%v, %v), got (%v, %v)",
					p1.X(), p1.Y(), messages[i], p2.X(), p2.Y(), expect.X(), expect.Y(), actual.X(), actual.Y(),
				)
			}
		}
	}
}

func TestTranslateOut(t *testing.T) {
	transforms := []Transform{
		MANDELBROT_BOUNDS,
		MANDELBROT_BOUNDS,
		MANDELBROT_BOUNDS,
		UHD,
		MINI,
	}
	cases := [][2]Point{
		[2]Point{MANDELBROT_LOWER_BOUND, ORIGIN_COMPLEX},
		[2]Point{MANDELBROT_UPPER_BOUND, MANDELBROT_RANGE},
		[2]Point{MANDELBROT_ARB_1, MA1_TRANSLATED_OUT},
		[2]Point{UHD_CENTER, UHD_CENTER},
		[2]Point{MINI_CENTER, MINI_CENTER},
	}
	messages := []string{
		"Mandelbrot Bounds",
		"Mandelbrot Bounds",
		"Mandelbrot Bounds",
	}
	for i, tc := range cases {
		p, expect := tc[0], tc[1]
		if actual := TranslateOut(transforms[i], p); actual != expect {
			t.Errorf("Translation out of %v expected (%v, %v), got (%v, %v)",
				messages[i],
				expect.X(), expect.Y(),
				actual.X(), actual.Y(),
			)
		}
	}
}

func TestTranslateInto(t *testing.T) {
	transforms := []Transform{
		MANDELBROT_BOUNDS,
		MANDELBROT_BOUNDS,
		MANDELBROT_BOUNDS,
		UHD,
		MINI,
	}
	cases := [][2]Point{
		[2]Point{ORIGIN_COMPLEX, MANDELBROT_LOWER_BOUND},
		[2]Point{MANDELBROT_RANGE, MANDELBROT_UPPER_BOUND},
		[2]Point{MA1_TRANSLATED_OUT, MANDELBROT_ARB_1},
		[2]Point{UHD_CENTER, UHD_CENTER},
		[2]Point{MINI_CENTER, MINI_CENTER},
	}
	messages := []string{
		"Mandelbrot Bounds",
		"Mandelbrot Bounds",
		"Mandelbrot Bounds",
		"Ultra HD",
		"Mini Render",
	}
	for i, tc := range cases {
		p, expect := tc[0], tc[1]
		if actual := TranslateInto(transforms[i], p); !PointApproxEq(actual, expect, .00001) {
			t.Errorf("Translation into %v expected (%v, %v), got (%v, %v)",
				messages[i],
				expect.X(), expect.Y(),
				actual.X(), actual.Y(),
			)
		}
	}
}

func TestTransformPoint(t *testing.T) {
	transforms := [][2]Transform{
		[2]Transform{MANDELBROT_BOUNDS, UHD},
		[2]Transform{MANDELBROT_BOUNDS, UHD},
		[2]Transform{MANDELBROT_BOUNDS, MINI},
		[2]Transform{UHD, MANDELBROT_BOUNDS},
		[2]Transform{MINI, MANDELBROT_BOUNDS},
	}
	cases := [][2]Point{
		[2]Point{ORIGIN_COMPLEX, UHD_CENTER},
		[2]Point{MANDELBROT_UPPER_BOUND, UHD_BOTTOM_RIGHT},
		[2]Point{MANDELBROT_LOWER_BOUND, ORIGIN_PIXEL},
		[2]Point{UHD_CENTER, ORIGIN_COMPLEX},
		[2]Point{MINI_CENTER, ORIGIN_COMPLEX},
	}
	messages := []string{
		"Mandelbrot Bounds -> Ultra HD",
		"Mandelbrot Bounds -> Ultra HD",
		"Mandelbrot Bounds -> Mini Render",
		"Ultra HD -> Mandelbrot Bounds",
		"Mini Render -> Mandelbrot Bounds",
	}
	for i, tc := range cases {
		p, expect := tc[0], tc[1]
		t1, t2 := transforms[i][0], transforms[i][1]
		if actual := TransformPoint(t1, t2, p); !PointApproxEq(actual, expect, .00001) {
			t.Errorf("Transform (%v, %v) %v expected (%v, %v), got (%v, %v)",
				p.X(), p.Y(),
				messages[i],
				expect.X(), expect.Y(),
				actual.X(), actual.Y(),
			)
		}
	}
}

func ApproxEq(v1 float64, v2 float64, tolerance float64) bool {
	return math.Abs(v1-v2) < tolerance
}

func PointApproxEq(p1 Point, p2 Point, tolerance float64) bool {
	return ApproxEq(p1.X(), p2.X(), tolerance) && ApproxEq(p1.Y(), p2.Y(), tolerance)
}
