package example

func _(xs []int) []int {
	const toRemove = 10
	var filtered []int
	filtered = xs[0:0]
	for i := int(0); i < len(xs); i++ {
		x := xs[i]
		if toRemove+1 != x {
			filtered = append(filtered, x)
		}
	}
	return (filtered)
}
