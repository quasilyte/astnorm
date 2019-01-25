#!/usr/bin/env bash

# Both functions implement strings.Repeat.
# Their source code differs.
cat mylib1/mylib1.go
cat mylib2/mylib2.go

# Normalized forms are more-or-less the same.
# Only variables names differ.
go-normalize mylib1/mylib1.go
go-normalize mylib2/mylib2.go

# Make normalized packages.
go-normalize mylib1/mylib1.go > mylib1n/mylib1.go
go-normalize mylib2/mylib2.go > mylib2n/mylib2.go

# Use grepfunc to create a pattern from Go code.
# With syntax patterns, we can now ignore variable
# names differences and find both functions
# by either of them.
grepfunc -input mylib1n/mylib1.go -pattern=makeString ./...
grepfunc -input mylib2n/mylib2.go -pattern=repeatString ./...
