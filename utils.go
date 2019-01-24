package astnorm

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/typep"
)

func isLiteralConst(info *types.Info, x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.Ident:
		// Not really literal consts, but they are
		// considered as such by many programmers.
		switch x.Name {
		case "nil", "true", "false":
			return true
		}
		return false
	case *ast.BasicLit:
		return true
	default:
		return false
	}
}

func isCommutative(info *types.Info, x *ast.BinaryExpr) bool {
	// TODO(quasilyte): make this list more or less complete.
	switch x.Op {
	case token.ADD:
		return !typep.HasStringProp(info.TypeOf(x))
	case token.EQL, token.NEQ:
		return true
	default:
		return false
	}
}
