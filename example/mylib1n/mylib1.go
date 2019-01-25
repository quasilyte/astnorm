package mylib1

import "strings"

func makeString(str string, num int) string {
	parts := make([]string, num)
	for i := range parts {
		parts[i] = str
	}
	return strings.Join(parts, "")
}
