package mylib1

import "strings"

func makeString(str string, num int) string {
	var parts = make([]string, num)
	for i := 0; i < len(parts); i++ {
		parts[i] = str
	}
	return strings.Join(parts, "")
}
