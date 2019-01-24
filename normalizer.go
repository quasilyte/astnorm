package astnorm

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"

	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/astequal"
	"github.com/go-toolsmith/astp"
	"github.com/go-toolsmith/typep"
)

type normalizer struct {
	cfg *Config
}

func newNormalizer(cfg *Config) *normalizer {
	return &normalizer{
		cfg: cfg,
	}
}

func (n *normalizer) foldConstexpr(x ast.Expr) ast.Expr {
	if astp.IsParenExpr(x) {
		return nil
	}
	tv := n.cfg.Info.Types[x]
	if tv.Value == nil {
		return nil
	}
	if lit, ok := x.(*ast.BasicLit); ok && lit.Kind == token.FLOAT {
		// Floats may require additional handling here.
		if tv.Value.Kind() == constant.Int {
			// Usually, this case means that value is 0.
			// But to be sure, keep this assertion here.
			v, exact := constant.Int64Val(tv.Value)
			if !exact || v != 0 {
				panic(fmt.Sprintf("unexpected value for float with kind=int"))
			}
			return &ast.BasicLit{
				Kind:  token.FLOAT,
				Value: "0.0",
			}
		}
	}
	return constValueNode(tv.Value)
}

func (n *normalizer) normalizeExpr(x ast.Expr) ast.Expr {
	if folded := n.foldConstexpr(x); folded != nil {
		return folded
	}

	switch x := x.(type) {
	case *ast.CallExpr:
		if typep.IsTypeExpr(n.cfg.Info, x.Fun) {
			return n.normalizeTypeConversion(x)
		} else {
			x.Fun = n.normalizeExpr(x.Fun)
		}
		x.Args = n.normalizeExprList(x.Args)
		return x
	case *ast.SliceExpr:
		return n.normalizeSliceExpr(x)
	case *ast.ParenExpr:
		return n.normalizeExpr(x.X)
	case *ast.BinaryExpr:
		return n.normalizeBinaryExpr(x)
	default:
		return x
	}
}

func (n *normalizer) normalizeTypeConversion(x *ast.CallExpr) ast.Expr {
	typeTo := n.cfg.Info.TypeOf(x)
	typeFrom := n.cfg.Info.TypeOf(x.Args[0])
	if types.Identical(typeTo, typeFrom) {
		return x.Args[0]
	}
	return x
}

func (n *normalizer) normalizeSliceExpr(x *ast.SliceExpr) *ast.SliceExpr {
	x.Low = n.normalizeExpr(x.Low)
	x.High = n.normalizeExpr(x.High)
	x.Max = n.normalizeExpr(x.Max)
	x.X = n.normalizeExpr(x.X)
	// Omit default low boundary.
	if astcast.ToBasicLit(x.Low).Value == "0" {
		x.Low = nil
	}
	// Omit default high boundary, but only if 3rd index is abscent.
	if x.Max == nil {
		lenCall := astcast.ToCallExpr(x.High)
		if astcast.ToIdent(lenCall.Fun).Name == "len" && astequal.Expr(lenCall.Args[0], x.X) {
			x.High = nil
		}
	}
	return x
}

func (n *normalizer) normalizeBinaryExpr(x *ast.BinaryExpr) ast.Expr {
	x.X = n.normalizeExpr(x.X)
	x.Y = n.normalizeExpr(x.Y)

	// TODO(quasilyte): implement this check in a proper way.
	// Also handle empty strings.
	switch {
	case isCommutative(n.cfg.Info, x) && astcast.ToBasicLit(x.X).Value == "0":
		return x.Y
	case astcast.ToBasicLit(x.Y).Value == "0":
		return x.X
	}

	if isCommutative(n.cfg.Info, x) {
		lhs := astcast.ToBinaryExpr(x.X)
		cv1 := constValueOf(n.cfg.Info, lhs.Y)
		cv2 := constValueOf(n.cfg.Info, x.Y)

		if cv1 != nil && cv2 != nil {
			cv := constant.BinaryOp(cv1, x.Op, cv2)
			x.X = lhs.X
			x.Y = constValueNode(cv)
			return n.normalizeExpr(x)
		}

		// Turn yoda expressions into the more conventional notation.
		// Put constant inside the expression after the non-constant part.
		if isLiteralConst(n.cfg.Info, x.X) && !isLiteralConst(n.cfg.Info, x.Y) {
			x.X, x.Y = x.Y, x.X
		}
	}

	return x
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
	case *ast.ReturnStmt:
		return n.normalizeReturnStmt(x)
	case *ast.DeclStmt:
		return n.normalizeDeclStmt(x)
	default:
		return x
	}
}

func (n *normalizer) normalizeReturnStmt(ret *ast.ReturnStmt) *ast.ReturnStmt {
	ret.Results = n.normalizeExprList(ret.Results)
	return ret
}

func (n *normalizer) normalizeBlockStmt(b *ast.BlockStmt) *ast.BlockStmt {
	list := b.List[:0]
	for _, x := range b.List {
		// Filter-out const decls.
		// We inline const values, so local const defs are
		// not needed to keep code valid.
		decl, ok := x.(*ast.DeclStmt)
		if ok && decl.Decl.(*ast.GenDecl).Tok == token.CONST {
			continue
		}
		list = append(list, n.normalizeStmt(x))
	}
	b.List = list

	n.normalizeValSwap(b)

	return b
}

func (n *normalizer) normalizeDeclStmt(stmt *ast.DeclStmt) ast.Stmt {
	decl := stmt.Decl.(*ast.GenDecl)
	if decl.Tok != token.VAR {
		return stmt
	}
	if len(decl.Specs) != 1 {
		return stmt
	}
	spec := decl.Specs[0].(*ast.ValueSpec)
	if len(spec.Names) != 1 {
		return stmt
	}

	switch {
	case len(spec.Values) == 1:
		// var x T = v
		return &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{spec.Names[0]},
			Rhs: []ast.Expr{spec.Values[0]},
		}
	case len(spec.Values) == 0 && spec.Type != nil:
		// var x T
		zv := zeroValueOf(n.cfg.Info.TypeOf(spec.Type))
		if zv == nil {
			return stmt
		}
		return &ast.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []ast.Expr{spec.Names[0]},
			Rhs: []ast.Expr{zv},
		}
	default:
		return stmt
	}
}

func (n *normalizer) normalizeValSwap(b *ast.BlockStmt) {
	// tmp := x
	// x = y
	// y = tmp
	//
	// =>
	//
	// x, y = y, x
	//
	// FIXME(quasilyte): if tmp is used somewhere outside of the value swap,
	// this transformation is illegal.

	for i := 0; i < len(b.List)-2; i++ {
		assignTmp := astcast.ToAssignStmt(b.List[i+0])
		assignX := astcast.ToAssignStmt(b.List[i+1])
		assignY := astcast.ToAssignStmt(b.List[i+2])
		if assignTmp.Tok != token.DEFINE {
			continue
		}
		if assignX.Tok != token.ASSIGN || assignY.Tok != token.ASSIGN {
			continue
		}
		if len(assignTmp.Lhs) != 1 || len(assignX.Lhs) != 1 || len(assignY.Lhs) != 1 {
			continue
		}
		tmp := astcast.ToIdent(assignTmp.Lhs[0])
		x := assignX.Lhs[0]
		y := assignY.Lhs[0]
		if !astequal.Expr(assignTmp.Rhs[0], x) {
			continue
		}
		if !astequal.Expr(assignX.Rhs[0], y) {
			continue
		}
		if !astequal.Expr(assignY.Rhs[0], tmp) {
			continue
		}

		b.List[i] = &ast.AssignStmt{
			Tok: token.ASSIGN,
			Lhs: []ast.Expr{x, y},
			Rhs: []ast.Expr{y, x},
		}
		b.List = append(b.List[:i+1], b.List[i+3:]...)
	}
}

func (n *normalizer) normalizeAssignStmt(assign *ast.AssignStmt) ast.Stmt {
	for i, lhs := range assign.Lhs {
		assign.Lhs[i] = n.normalizeExpr(lhs)
	}
	for i, rhs := range assign.Rhs {
		assign.Rhs[i] = n.normalizeExpr(rhs)
	}
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
		assign.Rhs[0] = rhs.Y
	}
	return assign
}
