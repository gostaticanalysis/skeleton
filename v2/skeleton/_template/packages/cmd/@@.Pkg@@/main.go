@@ if .Cmd -@@
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"@@.Pkg@@"
	"golang.org/x/tools/go/packages"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	@@.Pkg@@.Analyzer.Flags = flag.NewFlagSet(@@.Pkg@@.Analyzer.Name, flag.ExitOnError)
	@@.Pkg@@.Analyzer.Flags.Parse(os.Args[1:])

	if @@.Pkg@@.Analyzer.Flags.NArg() < 1 {
		return errors.New("patterns of packages must be specified")
	}

	pkgs, err := packages.Load(@@.Pkg@@.Analyzer.Config, @@.Pkg@@.Analyzer.Flags.Args()...)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		pass := &@@.Pkg@@.Pass{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Pkg:    pkg,
		}
		if err := @@.Pkg@@.Analyzer.Run(pass); err != nil {
			return err
		}
	}

	return nil
}
@@end@@
