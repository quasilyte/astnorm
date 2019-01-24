package astnorm

import (
	"go/ast"
	"strings"
	"testing"

	"github.com/go-toolsmith/astcast"
	"github.com/go-toolsmith/astequal"
	"github.com/go-toolsmith/astfmt"
	"golang.org/x/tools/go/packages"
)

func isTestCase(assign *ast.AssignStmt) bool {
	return len(assign.Lhs) == 2 &&
		len(assign.Rhs) == 2 &&
		astcast.ToIdent(assign.Lhs[0]).Name == "_" &&
		astcast.ToIdent(assign.Lhs[1]).Name == "_"
}

func TestNormalizeExpr(t *testing.T) {
	pkg := loadPackage(t, "./testdata/normalize_expr.go")
	funcs := collectFuncDecls(pkg)
	cfg := &Config{Info: pkg.TypesInfo}

	for _, fn := range funcs {
		for _, stmt := range fn.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if !ok || !isTestCase(assign) {
				continue
			}
			input := assign.Rhs[0]
			want := assign.Rhs[1]
			have := Expr(cfg, input)
			if !astequal.Expr(have, want) {
				pos := pkg.Fset.Position(assign.Pos())
				t.Errorf("%s:\nhave: %s\nwant: %s",
					pos, astfmt.Sprint(have), astfmt.Sprint(want))
			}
		}
	}
}

func TestNormalizeStmt(t *testing.T) {
	pkg := loadPackage(t, "./testdata/normalize_stmt.go")
	funcs := collectFuncDecls(pkg)
	cfg := &Config{Info: pkg.TypesInfo}

	for _, fn := range funcs {
		for _, stmt := range fn.Body.List {
			assign, ok := stmt.(*ast.AssignStmt)
			if !ok || !isTestCase(assign) {
				continue
			}
			input := assign.Rhs[0].(*ast.FuncLit).Body
			want := assign.Rhs[1].(*ast.FuncLit).Body
			have := Stmt(cfg, input)
			if !astequal.Stmt(have, want) {
				pos := pkg.Fset.Position(assign.Pos())
				t.Errorf("%s:\nhave: %s\nwant: %s",
					pos, astfmt.Sprint(have), astfmt.Sprint(want))
			}
		}
	}
}

func collectFuncDecls(pkg *packages.Package) []*ast.FuncDecl {
	var funcs []*ast.FuncDecl
	for _, f := range pkg.Syntax {
		for _, decl := range f.Decls {
			decl, ok := decl.(*ast.FuncDecl)
			if !ok || decl.Body == nil {
				continue
			}
			if !strings.HasSuffix(decl.Name.Name, "Test") {
				continue
			}
			funcs = append(funcs, decl)
		}
	}
	return funcs
}

func loadPackage(t *testing.T, path string) *packages.Package {
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		t.Fatalf("load %q: %v", path, err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		t.Fatalf("package %q loaded with errors", path)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package from %q path, got %d",
			path, len(pkgs))
	}
	return pkgs[0]
}
