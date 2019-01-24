package normalize_expr

func addInts(x, y int) int { return x + y }

func identityTest() {
	var x int
	type T int

	_, _ = x, x
	_, _ = 102, 102
	_, _ = x+1, x+1
}

func yodaTest() {
	var x int
	var s string
	var m map[int]int

	_, _ = 1+x, x+1
	_, _ = (nil != m), m != nil

	// Concat is not commutative.
	_, _ = "prefix"+s, "prefix"+s
}

// TODO(quasilyte): implement this after yoda tests.
/*
func foldArithTest() {
	var x int

	// Zeroes can be removed completely.
	_, _ = x+0, x
	_, _ = x+0+0, x
	_, _ = 0+x, x
	_, _ = 0+0+x, x
	_, _ = 0+x+0, x
	_, _ = 0+0+x+0, x
	_, _ = 0+x+0+0, x

	// For commutative ops fold it into a single op.
	_, _ = x+1, x+1
	_, _ = x+1+1, x+2
}
*/

func parenthesisRemovalTest() {
	var x int
	type T int

	_, _ = (x), x
	_, _ = ((*T)(&x)), (*T)(&x)
	_, _ = (addInts)(1, 2), addInts(1, 2)
	_, _ = addInts((1), (2)), addInts(1, 2)
}
