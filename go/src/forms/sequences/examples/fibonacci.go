package fibonacci

type Fibonacci struct {
	AInitial uint
	BInitial uint
	a        uint
	b        uint
}

func MakeFib(a uint, b uint) Fibonacci {
	return Fibonacci{a, b, a, b}
}

func (f *Fibonacci) Gen() chan uint {
	ch := make(chan uint)
	go func() {
		for true {
			val := f.Next()
			ch <- val
		}
	}()
	return ch
}

func (f *Fibonacci) Next() (val uint) {
	val = f.a
	f.b += f.a
	f.a = f.b - f.a
	return
}

var BASE_FIB Fibonacci = MakeFib(0, 1)
var MATURE_FIB Fibonacci = MakeFib(69, 420)
