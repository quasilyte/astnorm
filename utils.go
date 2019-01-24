package astnorm

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"

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
	switch cv.Kind() {
	case constant.Bool:
		if constant.BoolVal(cv) {
			return &ast.Ident{Name: "true"}
		}
		return &ast.Ident{Name: "false"}

	case constant.String:
		v := constant.StringVal(cv)
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: `"` + v + `"`,
		}

	case constant.Int:
		// For whatever reason, constant.Int can also
		// mean "float with 0 fractional part".
		v, exact := constant.Int64Val(cv)
		if !exact {
			return nil
		}
		return &ast.BasicLit{
			Kind:  token.INT,
			Value: fmt.Sprint(v),
		}

	case constant.Float:
		v, exact := constant.Float64Val(cv)
		if !exact {
			return nil
		}
		s := fmt.Sprint(v)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		return &ast.BasicLit{
			Kind:  token.FLOAT,
			Value: s,
		}

	default:
		panic("unimplemented")
	}
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
