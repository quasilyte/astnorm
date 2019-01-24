package astnorm

import (
	"testing"

	"github.com/go-toolsmith/astfmt"
	"github.com/go-toolsmith/strparse"
)

func TestNormalizeExpr(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// Expressions that are already in canonical form.
		// Can't be normalized further.
		{`x`, `x`},
		{`102`, `102`},
	}

	cfg := &Config{}
	for _, test := range tests {
		normalized := Expr(cfg, strparse.Expr(test.input))
		have := astfmt.Sprint(normalized)
		if have != test.want {
			t.Errorf("normalize(%q):\nhave: %s\nwant: %s",
				test.input, have, test.want)
		}
	}
}
