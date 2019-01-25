package normalize_stmt

func addInts(x, y int) int { return x + y }

func identityTest() {
	var x int

	_, _ = func() {
		x += 1
		x -= 1
	}, func() {
		x += 1
		x -= 1
	}
}

func incdecStmtTest() {
	var x int

	_, _ = func() {
		x++
		(x)++
	}, func() {
		x++
		x++
	}
}

func rangeStmtTest() {
	var xs []int

	_, _ = func() {
		for i := range xs[0:len(xs)] {
			_ = (i)
		}
	}, func() {
		for i := range xs {
			_ = i
		}
	}
}

func assignOpTest() {
	var x int

	_, _ = func() {
		x = x + 5
		x = x - 2
		x = x * 4
	}, func() {
		x += 5
		x -= 2
		x *= 4
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

	_, _ = func() {
		const n = 10
		x := 10
		_ = x != n+1
	}, func() {
		x := 10
		_ = x != 11
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

	_, _ = func() {
		var x int
		_ = x
	}, func() {
		x := 0
		_ = x
	}

	_, _ = func() {
		var xs [][]int
		var s string
		_ = xs
		_ = s
	}, func() {
		xs := [][]int(nil)
		s := ""
		_ = xs
		_ = s
	}

	_, _ = func() {
		var xs [8]string
		_ = xs
	}, func() {
		xs := [8]string{}
		_ = xs
	}

	_, _ = func() (float64, float32) {
		var x float64
		var y float32
		return x, y
	}, func() (float64, float32) {
		x := 0.0
		y := float32(0.0)
		return x, y
	}
}

func rangeLoopTest() {
	_, _ = func() {
		var xs []int
		for i := 0; i < len(xs); i++ {
			x := xs[i]
			_ = x
		}

		// Uses i+1 index.
		for i := 0; i < len(xs); i++ {
			x := xs[i+1]
			_ = x
		}

		// Doesn't assign elem.

		for i := 0; i < len(xs); i++ {
			_ = i
		}

		// TODO(quasilyte): more negative tests.
		// (Hint: use coverage to guide you, Luke!)
	}, func() {
		xs := []int(nil)
		for _, x := range xs {
			_ = x
		}

		for i := 0; i < len(xs); i++ {
			x := xs[i+1]
			_ = x
		}

		for i := 0; i < len(xs); i++ {
			_ = i
		}
	}

	_, _ = func() {
		var xs []int
		const toRemove = 10
		var filtered []int
		filtered = xs[0:0]
		for i := int(0); i < len(xs); i++ {
			x := xs[i]
			if toRemove+1 != x {
				filtered = append(filtered, x)
			}
		}
		_ = (filtered)
	}, func() {
		xs := []int(nil)
		filtered := []int(nil)
		filtered = xs[:0]
		for _, x := range xs {
			if x != 11 {
				filtered = append(filtered, x)
			}
		}
		_ = filtered
	}

	_, _ = func(xs []int) {
		for i := 0; i < len(xs); i++ {
			_ = xs[i]
		}
	}, func(xs []int) {
		for i := range xs {
			_ = xs[i]
		}
	}

	_, _ = func(xs []int) {
		for i := 0; i < len(xs); i++ {
			_ = xs[i+1]
		}
	}, func(xs []int) {
		for i := 0; i < len(xs); i++ {
			_ = xs[i+1]
		}
	}

	_, _ = func(xs []int) {
		for i := 0; i < len(xs); i++ {
			v := xs[i]
			_ = v
			_ = xs[i]
		}
	}, func(xs []int) {
		for i, v := range xs {
			_ = v
			_ = xs[i]
		}
	}

	_, _ = func(xs []int) {
		for i := 0; i < len(xs); i++ {
			v := xs[i]
			_ = v
			_ = xs[i]
			i++
		}
	}, func(xs []int) {
		for i := 0; i < len(xs); i++ {
			v := xs[i]
			_ = v
			_ = xs[i]
			i++
		}
	}
}

func exprStmtTest() {
	_, _ = func() {
		addInts((1), 0+0+0)
	}, func() {
		addInts(1, 0)
	}
}

func combinedTest() {
	var x int

	_, _ = func() {
		x = x + (2)
	}, func() {
		x += 2
	}
}
