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
	return x
}
