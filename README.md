[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/astnorm)](https://goreportcard.com/report/github.com/Quasilyte/astnorm)
[![GoDoc](https://godoc.org/github.com/Quasilyte/astnorm?status.svg)](https://godoc.org/github.com/Quasilyte/astnorm)
[![Build Status](https://travis-ci.org/Quasilyte/astnorm.svg?branch=master)](https://travis-ci.org/Quasilyte/astnorm)

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
filtered = xs[0:len(xs)]
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
filtered = xs[:]
for _, x := range xs {
        if x != 11 {
                filtered = append(filtered, x)
        }
}
return filtered
```
     
  </td></tr>
</table>
