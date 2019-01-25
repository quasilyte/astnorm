package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os/exec"
	"strings"

	"github.com/Quasilyte/astnorm/cmd/internal/loadfile"
	"github.com/go-toolsmith/astfmt"
	"github.com/go-toolsmith/typep"
)

func main() {
	log.SetFlags(0)

	input := flag.String("input", "",
		`input Go file with pattern function`)
	pattern := flag.String("pattern", "_pattern",
		`function to be interpreted as a pattern`)
	verbose := flag.Bool("v", false,
		`turn on debug output`)
	flag.Parse()

	if *input == "" {
		log.Panic("-input argument can't be empty")
	}
	targets := flag.Args()

	f, info, err := loadfile.ByPath(*input)
	if err != nil {
		log.Panicf("loadfile: %v", err)
	}

	var fndecl *ast.FuncDecl
	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.FuncDecl)
		if ok && decl.Name.Name == *pattern {
			fndecl = decl
			break
		}
	}
	if fndecl == nil {
		log.Panicf("found no `%s` func in %q", *pattern, targets[0])
	}
	if fndecl.Body == nil {
		log.Panic("external funcs are not supported")
	}
	pat := makeGogrepPattern(info, fndecl.Body)
	s := astfmt.Sprint(pat)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	if *verbose {
		fmt.Println(s)
	}

	gogrepArgs := []string{"-x", s}
	gogrepArgs = append(gogrepArgs, targets...)
	out, err := exec.Command("gogrep", gogrepArgs...).CombinedOutput()
	if err != nil {
		log.Panicf("run gogrep: %v: %s", err, out)
	}
	fmt.Print(string(out))
}

type visitor struct {
	info *types.Info
}

func (v *visitor) visitNode(x ast.Node) bool {
	switch x := x.(type) {
	case *ast.Ident:
		// Don't replace type names.
		if typep.IsTypeExpr(v.info, x) {
			return true
		}
		x.Name = "$" + x.Name
		return true
	case *ast.CallExpr:
		// Don't want to replace function names.
		for _, arg := range x.Args {
			ast.Inspect(arg, v.visitNode)
		}
		return false
	default:
		return true
	}
}

func makeGogrepPattern(info *types.Info, body *ast.BlockStmt) *ast.BlockStmt {
	v := &visitor{info: info}
	ast.Inspect(body, v.visitNode)
	return body
}
