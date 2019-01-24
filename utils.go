package astnorm

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"github.com/go-toolsmith/astcast"
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

func constValueNode(cv constant.Value) ast.Expr {
	var folded ast.Expr

	switch cv.Kind() {
	case constant.Bool:
		if constant.BoolVal(cv) {
			folded = &ast.Ident{Name: "true"}
		} else {
			folded = &ast.Ident{Name: "false"}
		}

	case constant.String:
		v := constant.StringVal(cv)
		folded = &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + v + `"`,
		}

	case constant.Int:
		v, exact := constant.Int64Val(cv)
		if !exact {
			return nil
		}
		folded = &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprint(v),
		}

	case constant.Float:
		v, exact := constant.Float64Val(cv)
		if !exact {
			return nil
		}
		folded = &ast.BasicLit{
			Kind:  token.FLOAT,
			Value: fmt.Sprint(v),
		}

	case constant.Complex:
		panic("unimplemented")
	}

	// TODO(Quasilyte): returned value now misses type info.
	// See also #1.
	return folded
}

func constValueOf(info *types.Info, x ast.Expr) constant.Value {
	if cv := info.Types[x].Value; cv != nil {
		return cv
	}
	lit := astcast.ToBasicLit(x)
	switch lit.Kind {
	case token.INT:
		v, err := strconv.ParseInt(lit.Value, 10, 64)
		if err != nil {
			return nil
		}
		return constant.MakeInt64(v)
	default:
		return nil
	}
}
