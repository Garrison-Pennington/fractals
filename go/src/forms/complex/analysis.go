package complex

type Window struct {
	Height int
	Width  int
	Array  [][]ComplexFractalValue
}

func (w Window) IntoArray(arr [][]ComplexFractalValue) Window {
	w.Array = arr
	return w
}

func (w Window) ArrayHeight() int {
	return len(w.Array)
}

func (w Window) ArrayWidth() int {
	return len(w.Array[0])
}

func (w Window) Slide(stride int) chan WindowView {
	ch := make(chan WindowView)
	arr := make([][]ComplexFractalValue, w.Height)
	go func() {
		defer close(ch)
		for x := 0; x+w.Width < w.ArrayWidth(); x += stride {
			for y := 0; y+w.Height < w.ArrayHeight(); y += stride {
				for r := range arr {
					arr[r] = w.Array[y+r][x : x+w.Width]
				}
				view := WindowView{&w, arr, x, y}
				ch <- view
			}
		}
	}()
	return ch
}

type WindowView struct {
	Window *Window
	Array  [][]ComplexFractalValue
	X      int
	Y      int
}

func (view WindowView) Average(score func(ComplexFractalValue) float64) (total float64) {
	count := 0
	for value := range GenFractalArray(view.Array) {
		total += score(value)
		count++
	}
	total /= float64(count)
	return
}

func (view WindowView) EdgeCount(maxSteps uint32) (count int) {
	arr := view.Array
	mapFn := BinStepScorer(maxSteps)
	for r := range arr {
		for c := range arr[r] {
			val := mapFn(arr[r][c])
			isBottom := r == len(arr)-1
			isRight := c == len(arr[r])-1
			if !isBottom && val != mapFn(arr[r+1][c]) {
				count++
			}
			if !isRight && val != mapFn(arr[r][c+1]) {
				count++
			}
		}
	}
	return
}

func (view WindowView) Bounds() Bounds {
	lower := MakeComplexPoint(view.Array[0][0].Initial)
	w, h := view.Window.Width, view.Window.Height
	upper := MakeComplexPoint(view.Array[w-1][h-1].Initial)
	return Bounds{lower, upper}
}

func StepScorer(maxSteps uint32) func(ComplexFractalValue) float64 {
	return func(val ComplexFractalValue) float64 {
		return float64(val.Steps) / float64(maxSteps)
	}
}

func BinStepScorer(maxSteps uint32) func(ComplexFractalValue) bool {
	return func(val ComplexFractalValue) bool {
		return val.Steps == maxSteps
	}
}
