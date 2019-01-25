[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/astnorm)](https://goreportcard.com/report/github.com/Quasilyte/astnorm)
[![GoDoc](https://godoc.org/github.com/Quasilyte/astnorm?status.svg)](https://godoc.org/github.com/Quasilyte/astnorm)
[![Build Status](https://travis-ci.org/Quasilyte/astnorm.svg?branch=master)](https://travis-ci.org/Quasilyte/astnorm)

![logo](/logo.jpg)

# astnorm

Go AST normalization experiment.

> THIS IS NOT A PROPER LIBRARY (yet?).<br>
> DO NOT USE.<br>
> It will probably be completely re-written before it becomes usable.

## Normalized code examples

1. Remove elements that are equal to `toRemove+1`.

<table>
  <tr>
    <th>Before</th>
    <th>After</th>
  </tr>
  
  <tr><td>
  
```go
const toRemove = 10
var filtered []int
filtered = xs[0:0]
for i := 0; i < len(xs); i++ {
        x := xs[i]
        if x != toRemove+1 {
                filtered = append(filtered, x)
        }
}
return (filtered)
```
  
  </td><td>
     
 ```go
filtered := []int{}
filtered = xs[:0]
for _, x := range xs {
        if x != 11 {
                filtered = append(filtered, x)
        }
}
return filtered
```
     
  </td></tr>
</table>

## Usage examples

* [cmd/go-normalize](/cmd/go-normalize): normalize given Go file
* [cmd/grepfunc](/cmd/grepfunc): turn Go code into a pattern for `gogrep` and run it

Potential workflow:

* Normalize the entire Go stdlib
* Then normalize your function
* Run `grepfunc` against normalized stdlib
* If function you implemented has implementation under the stdlib, you'll probably find it
