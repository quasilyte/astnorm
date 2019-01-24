package main

import (
	"flag"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/Quasilyte/astnorm"
	"github.com/Quasilyte/astnorm/cmd/internal/loadfile"
)

func main() {
	log.SetFlags(0)

	flag.Parse()

	targets := flag.Args()
	if len(targets) != 1 {
		log.Panicf("expected exactly 1 positional argument (input go file)")
	}

	// For now, handle only 1 input file case.
	// For simplicity reasons.
	f, info, err := loadfile.ByPath(targets[0])
	if err != nil {
		log.Panicf("loadfile: %v", err)
	}
	normalizationConfig := &astnorm.Config{
		Info: info,
	}
	f = normalizeFile(normalizationConfig, f)
	fset := token.NewFileSet()
	if err := printer.Fprint(os.Stdout, fset, f); err != nil {
		log.Panicf("print normalized file: %v", err)
	}
}

func normalizeFile(cfg *astnorm.Config, f *ast.File) *ast.File {
	// Strip comments.
	f.Doc = nil
	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			n.Doc = nil
		case *ast.GenDecl:
			n.Doc = nil
		case *ast.Field:
			n.Doc = nil
		case *ast.ImportSpec:
			n.Doc = nil
		case *ast.ValueSpec:
			n.Doc = nil
		case *ast.TypeSpec:
			n.Doc = nil
		default:
		}
		return true
	})
	f.Comments = nil

	for _, decl := range f.Decls {
		// TODO(quasilyte): could also normalize global vars,
		// consts and type defs, but funcs are OK for the POC.
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			decl.Body = astnorm.Block(cfg, decl.Body)
		}
	}
	return f
}
