package main

import "text/template"

var srcTempl = template.Must(template.New("pass.go").Parse(`package {{.Pkg}}

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "{{.Pkg}}",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "{{.Pkg}} is ..."

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.Ident:
			_ = n
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

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, {{.Pkg}}.Analyzer, "a")
}
`))

var adotgoTempl = template.Must(template.New("a.go").Parse(`package a

func main() {
	// want "pattern"
}
`))

var cmdMainTempl = template.Must(template.New("main.go").Parse(`package main

import (
	"{{.ImportPath}}"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main({{.Pkg}}.Analyzer) }
`))
