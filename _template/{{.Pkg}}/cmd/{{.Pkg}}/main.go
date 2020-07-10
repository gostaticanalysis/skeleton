{{ if .Cmd -}}
package main

import (
	"{{.ImportPath}}"
	"golang.org/x/tools/go/analysis/{{.Checker}}checker"
)

func main() { {{.Checker}}checker.Main({{.Pkg}}.Analyzer) }
{{end}}
