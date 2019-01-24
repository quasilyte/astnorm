package astnorm

import (
	"go/ast"
	"go/token"

	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/astequal"
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

func (n *normalizer) normalizeStmt(x ast.Stmt) ast.Stmt {
	switch x := x.(type) {
	case *ast.AssignStmt:
		return n.normalizeAssignStmt(x)
	case *ast.BlockStmt:
		return n.normalizeBlockStmt(x)
	default:
		return x
	}
}

func (n *normalizer) normalizeBlockStmt(b *ast.BlockStmt) *ast.BlockStmt {
	for i, x := range b.List {
		b.List[i] = n.normalizeStmt(x)
	}
	return b
}

func (n *normalizer) normalizeAssignStmt(assign *ast.AssignStmt) ast.Stmt {
	assign = n.normalizeAssignOp(assign)
	return assign
}

func (n *normalizer) normalizeAssignOp(assign *ast.AssignStmt) *ast.AssignStmt {
	if assign.Tok != token.ASSIGN || len(assign.Lhs) != 1 {
		return assign
	}
	rhs := astcast.ToBinaryExpr(assign.Rhs[0])
	if !astequal.Expr(assign.Lhs[0], rhs.X) {
		return assign
	}
	// TODO(quasilyte): add missing ops mapping.
	switch rhs.Op {
	case token.ADD:
		assign.Tok = token.ADD_ASSIGN
		assign.Rhs[0] = n.normalizeExpr(rhs.Y)
	}
	return assign
}
