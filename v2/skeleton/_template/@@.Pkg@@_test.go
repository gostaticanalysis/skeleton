@@ if (or (eq .Kind "inspect") (eq .Kind "ssa")) -@@
package @@.Pkg@@_test

import (
	"testing"

	"@@.ImportPath@@"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, @@.Pkg@@.Analyzer, "a")
}
@@ end -@@
@@ if eq .Kind "codegen" -@@
package @@.Pkg@@_test

import (
	"flag"
	"os"
	"testing"

	"@@.ImportPath@@"
	"github.com/gostaticanalysis/codegen/codegentest"
)

var flagUpdate bool

func TestMain(m *testing.M) {
	flag.BoolVar(&flagUpdate, "update", false, "update the golden files")
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerator(t *testing.T) {
	rs := codegentest.Run(t, codegentest.TestData(), @@.Pkg@@.Generator, "a")
	codegentest.Golden(t, rs, flagUpdate)
}
@@ end -@@
