package internal

import (
	"flag"
	"io"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

type Pass struct {
	*packages.Package
	SSA      *ssa.Program
	SrcFuncs []*ssa.Function
	Stdin    io.Reader
	Stdout   io.Writer
	Stderr   io.Writer
}

type Analyzer struct {
	Name           string
	Doc            string
	Flags          *flag.FlagSet
	Config         *packages.Config
	SSABuilderMode ssa.BuilderMode
	Run            func(pass *Pass) error
}
