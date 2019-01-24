[![Go Report Card](https://goreportcard.com/badge/github.com/Quasilyte/astnorm)](https://goreportcard.com/report/github.com/Quasilyte/astnorm)
[![GoDoc](https://godoc.org/github.com/Quasilyte/astnorm?status.svg)](https://godoc.org/github.com/Quasilyte/astnorm)
[![Build Status](https://travis-ci.org/Quasilyte/astnorm.svg?branch=master)](https://travis-ci.org/Quasilyte/astnorm)

# astnorm

Go AST normalization experiment.

> THIS IS NOT A PROPER LIBRARY (yet?).<br>
> DO NOT USE.<br>
> It will probably be completely re-written before it becomes usable.

## Normalized code examples

### Remove all zeros from xs, in-place

<table>
  <tr>
    <th>Before</th>
    <th>After</th>
  </tr>
  
  <tr><td>
  
```go
var filtered []int
filtered = xs[0:len(xs)]
for i := 0; i < len(xs); i++ {
        x := xs[i]
        if x != 0 {
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
        if x != 0 {
                filtered = append(filtered, x)
        }
}
return filtered
```
     
  </td></tr>
</table>
