package normalize_expr

func identity() {
	var x int
	type T int

	_, _ = x, x
	_, _ = 102, 102
}

func parenthesisRemoval() {
	var x int
	type T int

	_, _ = (x), x
	_, _ = ((*T)(&x)), (*T)(&x)
}
