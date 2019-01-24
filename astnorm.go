package astnorm

import (
	"go/ast"
	"go/types"
)

// Config carries information needed to properly normalize
// AST nodes as well as optional configuration values
// to control different aspects of the process.
type Config struct {
	Info *types.Info
}

// Expr returns normalized expression x.
// x may be mutated.
func Expr(cfg *Config, x ast.Expr) ast.Expr {
	return newNormalizer(cfg).normalizeExpr(x)
}

// Stmt returns normalized statement x.
// x may be mutated.
func Stmt(cfg *Config, x ast.Stmt) ast.Stmt {
	return newNormalizer(cfg).normalizeStmt(x)
}

// Block returns normalized block x.
// x may be mutated.
func Block(cfg *Config, x *ast.BlockStmt) *ast.BlockStmt {
	return newNormalizer(cfg).normalizeBlockStmt(x)
}
