package @@.Pkg@@_test

import (
	"flag"
	"os"
	"testing"

	"@@.Path@@"
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
