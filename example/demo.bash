#!/usr/bin/env bash

# Both packages have a func that implements strings.Repeat.
# Suppose we found 1 such place in our codebase and calls to it
# with strings.Repeat. But how maybe there are more such funcs?
#
# Note that their source code differs.
# One uses for loop, the other uses range form of loop.
# One uses const, the other don't.
cat mylib1/mylib1.go
## package mylib1
##
## import "strings"
##
## func makeString(str string, num int) string {
## 	var parts = make([]string, num)
## 	for i := 0; i < len(parts); i++ {
## 		parts[i] = str
## 	}
## 	return strings.Join(parts, "")
## }

cat mylib2/mylib2.go
## package mylib2
##
## import "strings"
##
## func repeatString(s string, n int) string {
## 	var pieces = make([]string, n)
## 	for i := range pieces {
## 		pieces[i] = s
## 	}
## 	const sep = ""
## 	return strings.Join(pieces, sep)
## }

# Normalized forms are more-or-less the same.
# Only variables names differ.
go-normalize mylib1/mylib1.go
## package mylib1
##
## import "strings"
##
## func makeString(str string, num int) string {
## 	parts := make([]string, num)
## 	for i := range parts {
## 		parts[i] = str
## 	}
## 	return strings.Join(parts, "")
## }

go-normalize mylib2/mylib2.go
## package mylib2
##
## import "strings"
##
## func repeatString(s string, n int) string {
## 	pieces := make([]string, n)
## 	for i := range pieces {
## 		pieces[i] = s
## 	}
## 	return strings.Join(pieces, "")
## }

# Make normalized packages.
go-normalize mylib1/mylib1.go > mylib1n/mylib1.go
go-normalize mylib2/mylib2.go > mylib2n/mylib2.go

# Use grepfunc to create a pattern from Go code.
# With syntax patterns, we can now ignore variable
# names differences and find both functions
# by either of them.
grepfunc -input mylib1n/mylib1.go -pattern=makeString ./...

## mylib1n/mylib1.go:6:2: parts := make([]string, num); for i := range parts { parts[i] = str; }; return strings.Join(parts, "")
## mylib2n/mylib2.go:6:2: pieces := make([]string, n); for i := range pieces { pieces[i] = s; }; return strings.Join(pieces, "")
grepfunc -input mylib2n/mylib2.go -pattern=repeatString ./...
## mylib1n/mylib1.go:6:2: parts := make([]string, num); for i := range parts { parts[i] = str; }; return strings.Join(parts, "")
## mylib2n/mylib2.go:6:2: pieces := make([]string, n); for i := range pieces { pieces[i] = s; }; return strings.Join(pieces, "")
