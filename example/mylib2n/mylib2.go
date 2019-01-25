package mylib2

import "strings"

func repeatString(s string, n int) string {
	pieces := make([]string, n)
	for i := range pieces {
		pieces[i] = s
	}
	return strings.Join(pieces, "")
}
