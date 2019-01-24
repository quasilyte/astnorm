package normalize_expr

func addInts(x, y int) int { return x + y }

func identityTest() {
	var x int
	type T int

	_, _ = x, x
	_, _ = 102, 102
}

func parenthesisRemovalTest() {
	var x int
	type T int

	_, _ = (x), x
	_, _ = ((*T)(&x)), (*T)(&x)
	_, _ = (addInts)(1, 2), addInts(1, 2)
	_, _ = addInts((1), (2)), addInts(1, 2)
}
