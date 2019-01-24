package astnorm

import (
	"go/ast"

	"github.com/go-toolsmith/typep"
)

type normalizer struct {
	cfg *Config
}

func newNormalizer(cfg *Config) *normalizer {
	return &normalizer{cfg: cfg}
}

func (n *normalizer) normalizeExpr(x ast.Expr) ast.Expr {
	switch x := x.(type) {
	case *ast.CallExpr:
		if !typep.IsTypeExpr(n.cfg.Info, x.Fun) {
			x.Fun = n.normalizeExpr(x.Fun)
		}
		x.Args = n.normalizeExprList(x.Args)
		return x
	case *ast.ParenExpr:
		return n.normalizeExpr(x.X)
	default:
		return x
	}
}

func (n *normalizer) normalizeExprList(xs []ast.Expr) []ast.Expr {
	for i, x := range xs {
		xs[i] = n.normalizeExpr(x)
	}
	return xs
}
