package @@.Pkg@@_test

import (
	"testing"

	"@@.Path@@"
	"github.com/gostaticanalysis/testutil"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	@@if .CopyParentGoMod -@@
	modfile := testutil.ModFile(t, ".", nil)
	testdata := testutil.WithModules(t, analysistest.TestData(), modfile)
	@@else -@@
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	@@end -@@
	analysistest.Run(t, testdata, @@.Pkg@@.Analyzer, "a")
}
