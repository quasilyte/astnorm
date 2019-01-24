package normalize_stmt

func identityTest() {
	var x int

	_, _ = func() {
		x += 1
	}, func() {
		x += 1
	}
}

func assignOpTest() {
	var x int

	_, _ = func() {
		x = x + 1
	}, func() {
		x += 1
	}
}

func combinedTest() {
	var x int

	_, _ = func() {
		x = x + (1)
	}, func() {
		x += 1
	}
}
