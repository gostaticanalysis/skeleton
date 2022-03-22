@@ if .Cmd -@@
package main

import (
	"@@.Path@@"
	"golang.org/x/tools/go/analysis/@@.Checker@@checker"
)

func main() { @@.Checker@@checker.Main(@@.Pkg@@.Analyzer) }
@@end@@
