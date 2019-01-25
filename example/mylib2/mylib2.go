package mylib2

import "strings"

func repeatString(s string, n int) string {
	var pieces = make([]string, n)
	for i := range pieces {
		pieces[i] = s
	}
	const sep = ""
	return strings.Join(pieces, sep)
}
