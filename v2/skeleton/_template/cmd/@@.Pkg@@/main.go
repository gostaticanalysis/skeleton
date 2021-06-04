@@ if .Cmd -@@
@@ if (or (eq .Kind "inspect") (eq .Kind "ssa")) -@@
package main

import (
	"@@.Path@@"
	"golang.org/x/tools/go/analysis/@@.Checker@@checker"
)

func main() { @@.Checker@@checker.Main(@@.Pkg@@.Analyzer) }
@@ end -@@
@@ if eq .Kind "codegen" -@@
package main

import (
	"@@.Path@@"
	"github.com/gostaticanalysis/codegen/@@.Checker@@generator"
)

func main() {
	@@.Checker@@generator.Main(@@.Pkg@@.Generator)
}
@@ end -@@
@@end@@
