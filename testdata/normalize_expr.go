package normalize_expr

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
}
