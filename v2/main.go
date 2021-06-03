package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"golang.org/x/mod/module"
)

//go:embed version.txt
var version string

func main() {

	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println("skeleton", version)
		os.Exit(0)
	}

	var info skeleton.Info
	parseFlag(&info)
	info.Path = flag.Arg(0)
	if module.CheckPath(info.Path) != nil {
		flag.Usage()
		os.Exit(1)
	}

	if info.Pkg == "" {
		info.Pkg = path.Base(info.Path)
	}

	if err := run(&info); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func parseFlag(info *skeleton.Info) {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "skeleton [-checker,-kind,-cmd,-plugin] example.com/path")
		flag.PrintDefaults()
	}
	flag.Var(&info.Checker, "checker", "[unit,single,multi]")	
	flag.Var(&info.Kind, "kind", "[inspect,ssa,codegen]")
	if info.Kind == "" {
		info.Kind = skeleton.KindInspect
	}
	flag.BoolVar(&info.Cmd, "cmd", false, "create main file")
	if info.Checker != "" {
		info.Checker = skeleton.CheckerUnit
		info.Cmd = true
	}
	flag.BoolVar(&info.Plugin, "plugin", false, "create golangci-lint plugin")
	flag.StringVar(&info.Pkg, "pkg", "", "package name")
	flag.Parse()
}

func run(info *skeleton.Info) error {
	fsys, err := new(skeleton.Generator).Run(context.Background(), info)
	if err != nil {
		return err
	}
	if err := skeleton.CreateDir(info.Pkg, fsys); err != nil {
		return err
	}
	return nil
}
