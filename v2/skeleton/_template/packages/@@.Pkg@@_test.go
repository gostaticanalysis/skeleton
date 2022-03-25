package @@.Pkg@@_test

import (
	"bytes"
	"flag"
	"path/filepath"
	"strings"
	"testing"

	"@@.Path@@"
	"@@.Path@@/internal"
	"github.com/tenntenn/golden"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func Test(t *testing.T) {
	pkgs := load(t, testdata(t), "a")
	for _, pkg := range pkgs {
		prog, funcs := buildssa(t, pkg)
		run(t, pkg, prog, funcs)
	}
}

func load(t *testing.T, testdata string, pkgname string) []*packages.Package {
	t.Helper()
	@@.Pkg@@.Analyzer.Config.Dir = filepath.Join(testdata, "src", pkgname)
	pkgs, err := packages.Load(@@.Pkg@@.Analyzer.Config, pkgname)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return pkgs
}

func buildssa(t *testing.T, pkg *packages.Package) (*ssa.Program, []*ssa.Function) {
	t.Helper()
	program, funcs, err := internal.BuildSSA(pkg, @@.Pkg@@.Analyzer.SSABuilderMode)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return program, funcs
}

func testdata(t *testing.T) string {
	t.Helper()
	dir, err := filepath.Abs("testdata")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	return dir
}

func run(t *testing.T, pkg *packages.Package, prog *ssa.Program, funcs []*ssa.Function) {
	var stdin, stdout, stderr bytes.Buffer
	pass := &internal.Pass{
		Stdin:    &stdin,
		Stdout:   &stdout,
		Stderr:   &stderr,
		Package:  pkg,
		SSA:      prog,
		SrcFuncs: funcs,
	}

	if err := @@.Pkg@@.Analyzer.Run(pass); err != nil {
		t.Error("unexpected error:", err)
	}

	pkgname := pkgname(pkg)

	if flagUpdate {
		golden.Update(t, testdata(t), pkgname+"-stdout", &stdout)
		golden.Update(t, testdata(t), pkgname+"-stderr", &stderr)
		return
	}

	if diff := golden.Diff(t, testdata(t), pkgname+"-stdout", &stdout); diff != "" {
		t.Errorf("stdout of analyzing %s:\n%s", pkgname, diff)
	}

	if diff := golden.Diff(t, testdata(t), pkgname+"-stderr", &stderr); diff != "" {
		t.Errorf("stderr of analyzing %s:\n%s", pkgname, diff)
	}
}

func pkgname(pkg *packages.Package) string {
	switch {
	case pkg.PkgPath != "":
		return strings.ReplaceAll(pkg.PkgPath, "/", "-")
	case pkg.Name != "":
		return pkg.Name
	default:
		return pkg.ID
	}
}
