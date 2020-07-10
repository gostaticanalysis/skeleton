{{ if .Cmd -}}
package main

import (
	"{{.ImportPath}}"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main({{.Pkg}}.Analyzer) }
{{end}}
