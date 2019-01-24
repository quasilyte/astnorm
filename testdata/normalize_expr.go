package normalize_expr

func addInts(x, y int) int { return x + y }

func identityTest() {
	var x int
	type T int

	_, _ = x, x
	_, _ = 102, 102
	_, _ = x+1, x+1
	_, _ = 0-x, 0-x
	_, _ = 1.1, 1.1
	_, _ = 12412.312, 12412.312
}

func defaultSlicingBoundsTest() {
	var xs []int
	var s string

	_, _ = xs[0:], xs[:]
	_, _ = (xs)[(0+0):], xs[:]
	_, _ = xs[0:len(xs)], xs[:]
	_, _ = (xs)[0:(len(xs))], xs[:]
	_, _ = xs[:0:0], xs[:0:0]

	_, _ = s[0:len(s)], s[:]
	_, _ = s[1:], s[1:]
}

func literalsTest() {
	// Convert any int numerical base into 10.
	_, _ = 0x0, 0
	_, _ = 0x1, 1
	_, _ = 04, 4
	_, _ = 010, 8

	// Represent floats in a consistent way.
	_, _ = 1.0, 1.0
	_, _ = 5.0, 5.0
	_, _ = 0.0, 0.0
	_, _ = .0, 0.0
	_, _ = 0., 0.0
	_, _ = 0.1e4, 1000.0
	_, _ = 00.0, 0.0
}

func conversionTest() {
	var x int

	// These alredy have proper type even without conversion.
	_, _ = int(1), 1
	_, _ = float64(40.1), 40.1
	_, _ = int(x), x
	_, _ = int(x+1), x+1

	// These require conversion.
	_, _ = int32(x), int32(x)
}

func yodaTest() {
	var x int
	var s string
	var m map[int]int

	_, _ = 1+x, x+1
	_, _ = (nil != m), m != nil

	// Concat is not commutative.
	_, _ = "prefix"+s, "prefix"+s
	// Other non-commutative ops.
	_, _ = 1-x, 1-x
	_, _ = 1000/x, 1000/x
}

// TODO(quasilyte): implement this after yoda tests.
func foldArithTest() {
	var x int

	// Const-only expressions are folded entirely.
	_, _ = 1+2+3, 6
	_, _ = 6-2, 4

	// Zeroes can be removed completely as well.
	_, _ = x+0, x
	_, _ = x+(0)+0, x
	_, _ = 0+x, x
	_, _ = 0+0+x, x
	_, _ = 0+x+(0), x
	_, _ = (0+0)+x+0, x
	_, _ = 0+x+0+0, x
	_, _ = x-0-0, x

	// For commutative ops fold it into a single op.
	_, _ = x+1, x+1
	_, _ = x+1+1, x+2
	_, _ = 1+x+1, x+2
	_, _ = 1+2+x+2+1, x+6
	_, _ = (1+2)+x+2+1, x+6
	_, _ = ((1 + (2)) + (x + 2) + 1), x+6
	_, _ = 0.2+0.1, 0.3

	_, _ = "a"+"b"+"c", "abc"
}

func parenthesisRemovalTest() {
	var x int
	type T int

	_, _ = (x), x
	_, _ = ((*T)(&x)), (*T)(&x)
	_, _ = (addInts)(1, 2), addInts(1, 2)
	_, _ = addInts((1), (2)), addInts(1, 2)
}
