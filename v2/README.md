# skeleton 

skeleton creates skeleton codes for a modularized static analysis tool with [x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package.

## x/tools/go/analysis pacakge

[x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package provides a type `analysis.Analyzer` which is unit of analyzers in modularized static analysis tool.

If you want to create new analyzer, you should provide a package variable which type is `*analysis.Analyzer`.
`skeleton` creates skeleton codes of the package and directories including test codes and main.go.

## Install

### Go version < 1.16

```
$ go get -u github.com/gostaticanalysis/skeleton/v2
```

### Go 1.16+

```
$ go install github.com/gostaticanalysis/skeleton/v2@latest
```

## How to use

### Create skeleton codes with module path

```
$ skeleton example.com/pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── go.mod
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

### Overwrite existing directory

If you want to overwrite without confirmation, you can run with `-overwrite` option.

```
$ skeleton -overwrite example.com/pkgname
```

If you run skeleton without `-overwrite` option, skeleton show optoins.
```
$ skeleton example.com/pkgname
pkgname already exist, remove?
[1] No(Exit)
[2] Remove and create new directory
[3] Overwrite existing files with confirmation
[4] Create new files only
```

### Create skeleton codes without cmd directory

```
$ skeleton -cmd=false example.com/pkgname
pkgname
├── go.mod
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

### Change the checker from unitchecker to singlechecker or multichecker

You can change the checker from unitchecker to singlechecker or multichecker.

```
$ skeleton -checker=single example.com/pkgname
$ cat cmd/pkgname/main.go
package main

import (
		"pkgname"
		"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(pkgname.Analyzer) }
```

### Create skeleton codes with plugin directory

```
$ skeleton -plugin example.com/pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── go.mod
├── pkgname.go
├── pkgname_test.go
├── plugin
│   └── main.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── go.mod
```

### Create skeleton codes of codegenerator

```
$ skeleton -type=codegen example.com/pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── go.mod
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            ├── a.go
            └── pkgname.golden
```

### Change type of skeleton code

skeleton accepts `-kind` option which indicates kind of skeleton code.

* `-kind=inspect`(default): generate skeleton code with `inspect.Analyzer`
* `-kind=ssa`: generate skeleton code with `buildssa.Analyzer`
* `-kind=codegen`: generate skeleton code of a code generator

## Build as a plugin for golangci-lint

`skeleton` generates plugin directory which has main.go.
The main.go can be built as a plugin for [golangci-lint](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint).

```
$ skeleton -plugin example.com/pkgname
$ go build -buildmode=plugin -o path_to_plugin_dir example.com/pkgname/plugin/pkgname
```

If you would like to specify flags for your plugin, you can put them via `ldflags` as below.

```
$ skeleton -plugin example.com/pkgname
$ go build -buildmode=plugin -ldflags "-X 'main.flags=-funcs log.Fatal'" -o path_to_plugin_dir example.com/pkgname/plugin/pkgname
```

### Without go.mod file

If you give `-gomod=false` flag to skeleton, skeleton does not create a go.mod file.
