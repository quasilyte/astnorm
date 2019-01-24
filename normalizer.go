package astnorm

import (
	"go/ast"
)

type normalizer struct {
	cfg *Config
}

func newNormalizer(cfg *Config) *normalizer {
	return &normalizer{cfg: cfg}
}

func (n *normalizer) normalizeExpr(x ast.Expr) ast.Expr {
	switch x := x.(type) {
	case *ast.ParenExpr:
		return n.normalizeExpr(x.X)
	default:
		return x
	}
}
