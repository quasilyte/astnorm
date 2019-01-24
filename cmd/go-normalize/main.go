package main

import (
	"flag"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"

	"github.com/Quasilyte/astnorm"
	"golang.org/x/tools/go/packages"
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
	// All those checks are to give more clear errors to the user.
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(cfg, targets...)
	if err != nil {
		log.Panicf("load: %v", err)
	}
	if errCount := packages.PrintErrors(pkgs); errCount != 0 {
		log.Panicf("%d errors during package loading", errCount)
	}
	if len(pkgs) != 1 {
		log.Panicf("loaded %d packages, expected only 1", len(pkgs))
	}
	pkg := pkgs[0]
	if len(pkg.Syntax) != 1 {
		log.Panicf("loaded package has %d files, expected only 1",
			len(pkg.Syntax))
	}

	normalizationConfig := &astnorm.Config{
		Info: pkg.TypesInfo,
	}
	f := normalizeFile(normalizationConfig, pkg.Syntax[0])
	fset := token.NewFileSet()
	if err := printer.Fprint(os.Stdout, fset, f); err != nil {
		log.Panicf("print normalized file: %v", err)
	}
}

func normalizeFile(cfg *astnorm.Config, f *ast.File) *ast.File {
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
