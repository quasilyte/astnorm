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

func valueSwapTest() {
	var x, y int

	_, _ = func() {
		tmp := (x)
		x = y
		y = tmp
	}, func() {
		x, y = y, x
	}

	_, _ = func() {
		tmp1 := x
		x = y
		y = tmp1

		tmp2 := y
		y = x
		x = tmp2
	}, func() {
		x, y = y, x
		y, x = x, y
	}

}

func removeConstDeclsTest() {
	_, _ = func() {
		const n = 10
		_ = n + n
	}, func() {
		_ = 20
	}
}

func rewriteVarSpecTest() {
	_, _ = func() {
		var x = 10
		var y float32 = float32(x)
		_ = x
		_ = y
	}, func() {
		x := 10
		y := float32(x)
		_ = x
		_ = y
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
