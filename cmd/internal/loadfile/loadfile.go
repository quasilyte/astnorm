package loadfile

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func ByPath(path string) (*ast.File, *types.Info, error) {
	cfg := &packages.Config{Mode: packages.LoadSyntax}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, nil, fmt.Errorf("load: %v", err)
	}
	if errCount := packages.PrintErrors(pkgs); errCount != 0 {
		return nil, nil, fmt.Errorf("%d errors during package loading", errCount)
	}
	if len(pkgs) != 1 {
		return nil, nil, fmt.Errorf("loaded %d packages, expected only 1", len(pkgs))
	}
	pkg := pkgs[0]
	if len(pkg.Syntax) != 1 {
		err := fmt.Errorf("loaded package has %d files, expected only 1",
			len(pkg.Syntax))
		return nil, nil, err
	}
	return pkg.Syntax[0], pkg.TypesInfo, nil
}
