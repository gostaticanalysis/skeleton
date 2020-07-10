# skeleton 

skeleton creates skeleton codes for a modularized static analysis tool with [x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package.

## x/tools/go/analysis pacakge

[x/tools/go/analysis](https://golang.org/x/tools/go/analysis) package provides a type `analysis.Analyzer` which is unit of analyzers in modularized static analysis tool.

If you want to create new analyzer, you should provide a package variable which type is `*analysis.Analyzer`.
`skeleton` creates skeleton codes of the package and directories including test codes and main.go.

## Insall

```
$ go get -u github.com/gostaticanalysis/skeleton
```

## How to use

### Create skeleton codes in GOPATH

```
$ skeleton pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── plugin
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Create skeleton codes with import path

```
$ skeleton -path="github.com/gostaticanalysis/pkgname"
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── plugin
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Overwrite existing directory

If you want to overwrite without confirmation, you can run with `-overwrite` option.

```
$ skeleton -overwrite pkgname
```

If you run skeleton without `-overwrite` option, skeleton show optoins.
```
$ skeleton pkgname
pkgname already exist, remove?
[1] No(Exit)
[2] Remove and create new directory
[3] Overwrite existing files with confirmation
[4] Create new files only
```

### Create skeleton codes without cmd directory

```
$ skeleton -cmd=false pkgname
pkgname
├── plugin
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

### Change the checker from unitchecker to singlechecker or multichecker

You can change the checker from unitchecker to singlechecker or multichecker.

```
$ skeleton -checker=single pkgname
$ cat cmd/pkgname/main.go                                                                    [~/Desktop/hogera]
package main

import (
		"pkgname"
		"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(pkgname.Analyzer) }
```

### Create skeleton codes without plugin directory

```
$ skeleton -plugin=false pkgname
pkgname
├── cmd
│   └── pkgname
│       └── main.go
├── pkgname.go
├── pkgname_test.go
└── testdata
    └── src
        └── a
            └── a.go
```

## Build as a plugin for golangci-lint

`skeleton` generates plugin directory which has main.go.
The main.go can be built as a plugin for [golangci-lint](https://golangci-lint.run/contributing/new-linters/#how-to-add-a-private-linter-to-golangci-lint).

```
$ skeleton pkgname
$ go build -buildmode=plugin -o path_to_plugin_dir importpath
```

If you would like to specify flags for your plugin, you can put them via `ldflags` as below.

```
$ skeleton pkgname
$ go build -buildmode=plugin -ldflags "-X 'main.flags=-funcs log.Fatal'" -o path_to_plugin_dir github.com/gostaticanalysis/called/plugin/called
```
