@@ if .Cmd -@@
@@ if (or (eq .Type "inspect") (eq .Type "ssa")) -@@
package main

import (
	"@@.ImportPath@@"
	"golang.org/x/tools/go/analysis/@@.Checker@@checker"
)

func main() { @@.Checker@@checker.Main(@@.Pkg@@.Analyzer) }
@@ end -@@
@@ if eq .Type "codegen" -@@
package main

import (
	"@@.ImportPath@@"
	"github.com/gostaticanalysis/codegen/@@.Checker@@generator"
)

func main() {
	@@.Checker@@generator.Main(@@.Pkg@@.Generator)
}
@@ end -@@
@@end@@
