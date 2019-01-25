#!/usr/bin/env bash

# Both packages have a func that implements strings.Repeat.
# Suppose we found 1 such place in our codebase and calls to it
# with strings.Repeat. But how maybe there are more such funcs?

cat mylib1/mylib1.go
go-normalize mylib1/mylib1.go > mylib1n/mylib1.go
diff mylib1/mylib1.go mylib1n/mylib1.go

cat mylib2/mylib2.go
go-normalize mylib2/mylib2.go > mylib2n/mylib2.go
diff mylib2/mylib2.go mylib2n/mylib2.go

# Use grepfunc to create a pattern from Go code.
# With syntax patterns, we can now ignore variable
# names differences and find both functions
# by either of them.
grepfunc -input mylib1n/mylib1.go -pattern=makeString ./...
grepfunc -v -input mylib2n/mylib2.go -pattern=repeatString ./...

# For more examples see consult README file.
