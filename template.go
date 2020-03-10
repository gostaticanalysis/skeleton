package main

import "text/template"

var srcTempl = template.Must(template.New("pass.go").Parse(`package {{.Pkg}}

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "{{.Pkg}} is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "{{.Pkg}}",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			if n.Name == "gopher" {
				pass.Reportf(n.Pos(), "identifyer is gopher")
			}
		}
	})

	return nil, nil
}
`))

var testTempl = template.Must(template.New("pass_test.go").Parse(`package {{.Pkg}}_test

import (
	"testing"

	"{{.ImportPath}}"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, {{.Pkg}}.Analyzer, "a")
}
`))

var adotgoTempl = template.Must(template.New("a.go").Parse(`package a

func f() {
	// The pattern can be written in regular expression.
	var gopher int // want "pattern"
	print(gopher)  // want "identifyer is gopher"
}
`))

var cmdMainTempl = template.Must(template.New("main.go").Parse(`package main

import (
	"{{.ImportPath}}"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main({{.Pkg}}.Analyzer) }
`))
