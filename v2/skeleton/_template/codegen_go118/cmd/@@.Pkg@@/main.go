@@ if .Cmd -@@
package main

import (
	"@@.Path@@"
	"github.com/gostaticanalysis/codegen/@@.Checker@@generator"
)

func main() {
	@@.Checker@@generator.Main(@@.Pkg@@.Generator)
}
@@end@@
