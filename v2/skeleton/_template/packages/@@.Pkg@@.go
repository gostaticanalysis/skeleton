package @@.Pkg@@

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"@@.Path@@/internal"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var Analyzer = &internal.Analyzer{
	Name: "@@.Pkg@@",
	Doc:  "@@.Pkg@@ is ...",
	Config: &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes |
			packages.NeedSyntax | packages.NeedTypesInfo |
			packages.NeedModule,
	},
	SSABuilderMode: 0,
	Run:            run,
}

func run(pass *internal.Pass) error {
	inspect := inspector.New(pass.Syntax)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pos := pass.Fset.Position(n.Pos())
				fname := pos.Filename
				if pass.Module != nil {
					var err error
					fname, err = filepath.Rel(pass.Module.Dir, fname)
					if err != nil {
						return
					}
				}
				fmt.Fprintf(pass.Stdout, "%s:%d:%d identifier is gopher\n", fname, pos.Line, pos.Column)
			}
		}
	})

	// See: golang.org/x/tools/go/ssa
	for _, f := range pass.SrcFuncs {
		fmt.Fprintln(pass.Stdout, f)
		for _, b := range f.Blocks {
			fmt.Fprintf(pass.Stdout, "\tBlock %d\n", b.Index)
			for _, instr := range b.Instrs {
				fmt.Fprintf(pass.Stdout, "\t\t%[1]T\t%[1]v\n", instr)
				for _, v := range instr.Operands(nil) {
					if v != nil {
						fmt.Fprintf(pass.Stdout, "\t\t\t%[1]T\t%[1]v\n", *v)
					}
               			}
			}
		}
	}

	return nil
}
