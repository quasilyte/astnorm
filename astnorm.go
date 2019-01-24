package astnorm

import (
	"go/ast"
	"go/types"

	"github.com/go-toolsmith/astcopy"
)

type Config struct {
	Info *types.Info
}

func Expr(cfg *Config, x ast.Expr) ast.Expr {
	copied := astcopy.Expr(x)
	return newNormalizer(cfg).normalizeExpr(copied)
}
