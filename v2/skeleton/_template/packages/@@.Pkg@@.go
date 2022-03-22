package @@.Pkg@@

import (
	"flag"
	"fmt"
	"go/ast"
	"io"
	"path/filepath"

	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

type Pass struct {
	Pkg    *packages.Package
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

var Analyzer = struct {
	Name   string
	Doc    string
	Flags  *flag.FlagSet
	Config *packages.Config
	Run    func(pass *Pass) error
}{
	Name: "@@.Pkg@@",
	Doc:  "@@.Pkg@@ is ...",
	Config: &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedModule,
	},
	Run: run,
}

func run(pass *Pass) error {
	inspect := inspector.New(pass.Pkg.Syntax)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pos := pass.Pkg.Fset.Position(n.Pos())
				fname, err := filepath.Rel(pass.Pkg.Module.Dir, pos.Filename)
				if err != nil {
					return
				}
				fmt.Fprintf(pass.Stdout, "%s:%d:%d identifier is gopher\n", fname, pos.Line, pos.Column)
			}
		}
	})

	return nil
}
