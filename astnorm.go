package astnorm

import (
	"go/ast"
	"go/types"
)

type Config struct {
	Info *types.Info
}

func Expr(cfg *Config, x ast.Expr) ast.Expr {
	return newNormalizer(cfg).normalizeExpr(x)
}

func Stmt(cfg *Config, x ast.Stmt) ast.Stmt {
	return newNormalizer(cfg).normalizeStmt(x)
}

func Block(cfg *Config, x *ast.BlockStmt) *ast.BlockStmt {
	return newNormalizer(cfg).normalizeBlockStmt(x)
}
