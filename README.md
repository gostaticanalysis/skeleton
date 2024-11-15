[日本語版](./README_ja.md)

# skeleton

skeleton is skeleton codes generator for Go's static analysis tools. skeleton makes easy to develop static analysis tools with [x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package and [x/tools/go/packages](https://golang.org/x/tools/go/packages) package.

## x/tools/go/analysis package

[x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package is for modularizing static analysis tools. x/tools/go/analysis package provides [analysis.Analyzer](https://golang.org/x/tools/go/analysis/#Analyzer) type which represents a unit of modularized static analysis tool.

`x/tools/go/analysis` package also provides common works of a static analysis tool. Just run the `skeleton mylinter` command, skeleton generates an `*analyzer.Analyzer` type initialization code, a test code, and a `main.go` for an executable which may be run with the `go vet` command.

The following blog helps to learn about the skeleton.

* [Go static analysis starting with skeleton](https://engineering.mercari.com/blog/entry/20220406-eea588f493/) (Japanese)

The following slides describes details of Go's static analysis including the `x/tools/go/analysis` package.

* [A complete introduction of the programming language Go, Chapter 14: Static Analysis and Code Generation](http://tenn.in/analysis) (Japanese)

## Installation

```
$ go install github.com/gostaticanalysis/skeleton/v2@latest
```

## How to use

### Create a skeleton code with a module path

skeleton receives a module path and generates a skeleton code with the module path. All generated codes are located in a directory which name is the last element of the module path.

When you run skeleton with `example.com/mylinter` as a module path, skeleton generates the following files.

```
$ skeleton example.com/mylinter
mylinter
├── cmd
│   └── mylinter
│       └── main.go
├── go.mod
├── mylinter.go
├── mylinter_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

#### Analyzer

A static analysis tool which developed with `x/tools/go/analysis`, is represented by value of `*analysis.Analyzer` type. In the mylinter case, the value is defined in `mylinter.go` as a variable which name is `Analyzer`.

The generated code provides toy implement with `inspect.Analyzer`. It finds identifiers which name are `gopher`.

```go
package mylinter

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "mylinter is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "mylinter",
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
				pass.Reportf(n.Pos(), "identifier is gopher")
			}
		}
	})

	return nil, nil
}
```

#### Test codes

skeleton also generates test codes. `x/tools/go/analysis` package provides a testing library in `analysistest` sub package. `analysistest.Run` runs tests with source codes in `testdata/src` directory. The second parameter is a path for `testdata` directory. The third parameter is test target analyzer and  remains are packages which are used in tests.

```go
package mylinter_test

import (
	"testing"

	"github.com/gostaticanalysis/example.com/mylinter"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, mylinter.Analyzer, "a")
}
```

In the mylinter case, the test uses `testdata/src/a/a.go` file as a test data. `mylinter.Analyzer` finds `gopher` identifiers in the source code and report them. In the test side, expected reports are described in comments. The comments must be start with `want` and a reporting message follows. The reporting message is represented by a regular expression. When the analyzer reports unexpected diagnostics or does not report expected diagnostics, the test will be failed.

```go
package a

func f() {
	// The pattern can be written in regular expression.
	var gopher int // want "pattern"
	print(gopher)  // want "identifier is gopher"
}
```

When you run `go mod tidy` and `go test`, the test will be failed because the analyzer does not report a diagnostic with "pattern".

```
$ go mod tidy
go: finding module for package golang.org/x/tools/go/analysis
go: finding module for package github.com/gostaticanalysis/testutil
go: finding module for package golang.org/x/tools/go/analysis/passes/inspect
go: finding module for package golang.org/x/tools/go/analysis/unitchecker
go: finding module for package golang.org/x/tools/go/ast/inspector
go: finding module for package golang.org/x/tools/go/analysis/analysistest
go: found golang.org/x/tools/go/analysis in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/analysis/passes/inspect in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/ast/inspector in golang.org/x/tools v0.1.10
go: found golang.org/x/tools/go/analysis/unitchecker in golang.org/x/tools v0.1.10
go: found github.com/gostaticanalysis/testutil in github.com/gostaticanalysis/testutil v0.4.0
go: found golang.org/x/tools/go/analysis/analysistest in golang.org/x/tools v0.1.10

$ go test
--- FAIL: TestAnalyzer (0.06s)
    analysistest.go:454: a/a.go:5:6: diagnostic "identifier is gopher" does not match pattern `pattern`
    analysistest.go:518: a/a.go:5: no diagnostic was reported matching `pattern`
FAIL
exit status 1
FAIL	github.com/gostaticanalysis/example.com/mylinter	1.270s
```

#### Executable file

skeleton generates `main.go` in `cmd` directory. When you build it and generate an executable file, the executable file must be run via `go vet` command such as the following. `-vettool` flag for `go vet` command specifies an absoluted path for an executable file of own static analysis tool.

```
$ go vet -vettool=`which mylinter` ./...
```

### Overwrite a directory

If the directory already exists, skeleton gives you with following options.

```
$ skeleton example.com/mylinter
mylinter already exists, overwrite?
[1] No (Exit)
[2] Remove and create new directory
[3] Overwrite existing files with confirmation
[4] Create new files only
```

### Without cmd directory

If you don't need `cmd` directory, you can set `false` to `-cmd` flag.

```
$ skeleton -cmd=false example.com/mylinter
mylinter
├── go.mod
mylinter.go
├── mylinter_test.go
└─ testdata
    └─ testdata
        testdata └── src
            Go.mod
            go.mod
```

### Without go.mod file

skeleton generates a `go.mod` file by default. When you would like to use skeleton in a directory which is already under Go Modules management, you can set `false` to `-gomod` option as following.

```
$ skeleton -gomod=false example.com/mylinter
mylinter
├── cmd
│└── mylinter
└─ main.go
├── mylinter.go
mylinter_test.go
└─ testdata
    testdata └── src
        testdata └─ a
            Go.mod
            go.mod
```

### SKELETON_PREFIX environment variable

When `SKELETON_PREFIX` environment variable is set, skeleton puts it as a prefix to a module path.

```
$ SKELETON_PREFIX=example.com skeleton mylinter
$ head -1 mylinter/go.mod
module example.com/mylinter
```

It is useful with [direnv](https://github.com/direnv/direnv) such as following.

```
$ cat ~/repos/gostaticanalysis/.envrc
export SKELETON_PREFIX=github.com/gostaticanalysis
```

If `SKELETON_PREFIX` environment variable is specified but the `-gomod` flag is `false`, skeleton prioritizes `-gomod` flag.

### singlechecker and multichecker

skeleton uses `unitchecker` package in `main.go` by default. You can change it to `singlechecker` package or `multichecker` package by specifying the `-checker` flag.

`singlechecker` package runs a single analyzer and `multichecker` package runs multiple analyzers. These packages does not need `go vet` command to run.

The following is an example of using `singlechecker` package.

```
$ skeleton -checker=single example.com/mylinter
$ cat cmd/mylinter/main.go
package main

import (
		"mylinter"
		"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(mylinter.Analyzer) }
```

Using `singlechecker` package or `multichecker` package seems easy way. But when you use them, you cannot receive benefit of using `go vet`. If you don't have particular reason of using `singlechecker` package or `multichecker` package, you should use `unitchecker`. It means you should not use `-checker` flag in most cases.

### Kinds of skeleton code

skeleton can change kind of skeleton code by using `-kind` flag.

* `-kind=inspect` (default): using `inspect.Analyzer`
* `-kind=ssa`: using the static single assignment (SSA, Static Single Assignment) form generated by `buildssa.Analyzer`
* `-kind=codegen`: code generator.
* `-kind=packages`: using `x/tools/go/packages` package

### Create code generator

When you gives `codegen` to `-kind` flag, skeleton generates skeleton code of code generation tool with [gostaticanalysis/codegen](https://pkg.go.dev/github.com/gostaticanalysis/codegen) package. 

```
$ skeleton -kind=codegen example.com/mycodegen
mycodegen
├── cmd
│   └── mycodegen
│       └── main.go
├── go.mod
├── mycodegen.go
├── mycodegen_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            ├── go.mod
            └── mycodegen.golden
```

`gostaticanalysis/codegen` package is an experimental, please be careful.

### golangci-lint plugin

skeleton generates codes that can be used as a plugin of [golangci-lint](https://github.com/golangci/golangci-lint) by specifying `-plugin` flag.

```
$ skeleton -plugin example.com/mylinter
mylinter
├── cmd
│   └── mylinter
│       └── main.go
├── go.mod
├── mylinter.go
├── mylinter_test.go
├── plugin
│   └── main.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

You can see [the documentation](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint). 

```
$ skeleton -plugin example.com/mylinter
$ go build -buildmode=plugin -o path_to_plugin_dir example.com/mylinter/plugin/mylinter
```

skeleton provides a way which can specify flags to your plugin with `-ldflags`. If you would like to know the details of it, please read the generated skeleton code.

```
$ skeleton -plugin example.com/mylinter
$ go build -buildmode=plugin -ldflags "-X 'main.flags=-funcs log.Fatal'" -o path_to_plugin_dir example.com/mylinter/plugin/mylinter
```

golangci-lint is built with `CGO_ENABLED=0` by default. So you should rebuilt with `CGO_ENABLED=1` because plugin package in the standard library uses CGO. And you should same version of modules with golangci-lint such as `golang.org/x/tools/go` module. The plugin system for golangci-lint is not recommended way.
