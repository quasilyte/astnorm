#!/usr/bin/env bash

# Both packages have a func that implements strings.Repeat.
# Suppose we found 1 such place in our codebase and calls to it
# with strings.Repeat. But how maybe there are more such funcs?
#
# Note that their source code differs.
# One uses for loop, the other uses range form of loop.
# One uses const, the other don't.
cat mylib1/mylib1.go
cat mylib2/mylib2.go

# Normalized forms are more-or-less the same.
# Only variables names differ.
go-normalize mylib1/mylib1.go

go-normalize mylib2/mylib2.go

# Make normalized packages.
go-normalize mylib1/mylib1.go > mylib1n/mylib1.go
go-normalize mylib2/mylib2.go > mylib2n/mylib2.go

# Compare mylib1.go before and after normalization.
vimdiff mylib1/mylib1.go mylib1n/mylib1.go

# Compare normalized mylib1.go and normalized mylib2.go.
vimdiff mylib1n/mylib1.go mylib2n/mylib2.go

# Use grepfunc to create a pattern from Go code.
# With syntax patterns, we can now ignore variable
# names differences and find both functions
# by either of them.
grepfunc -input mylib1n/mylib1.go -pattern=makeString ./...
grepfunc -input mylib2n/mylib2.go -pattern=repeatString ./...
