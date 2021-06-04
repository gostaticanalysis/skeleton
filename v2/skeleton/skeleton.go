package skeleton

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/mod/module"
)

const (
	ExitSuccess = 0
	ExitError   = 1
)

type Skeleton struct {
	Dir       string
	Output    io.Writer
	ErrOutput io.Writer
	Input     io.Reader
}

func Main(version string, args []string) int {
	s := &Skeleton{
		Dir:       ".",
		Output:    os.Stdout,
		ErrOutput: os.Stderr,
		Input:     os.Stdin,
	}
	return s.Run(version, args)
}

func (s *Skeleton) Run(version string, args []string) int {
	if len(args) > 0 && args[0] == "-v" {
		fmt.Fprintln(s.Output, "skeleton", version)
		return ExitSuccess
	}

	var info Info
	args, err := s.parseFlag(args, &info)
	if err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return ExitError
	}

	if len(args) <= 0 || module.CheckPath(args[0]) != nil {
		flag.Usage()
		return ExitError
	}
	info.Path = args[0]

	if info.Pkg == "" {
		info.Pkg = path.Base(info.Path)
	}

	if err := s.run(&info); err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return ExitError
	}

	return ExitSuccess
}

func (s *Skeleton) parseFlag(args []string, info *Info) ([]string, error) {
	flags := flag.NewFlagSet("skeleton", flag.ContinueOnError)
	flags.SetOutput(s.ErrOutput)
	flags.Usage = func() {
		fmt.Fprintln(s.ErrOutput, "skeleton [-checker,-kind,-cmd,-plugin] example.com/path")
		flags.PrintDefaults()
	}
	flags.Var(&info.Checker, "checker", "[unit,single,multi]")
	if info.Checker == "" {
		info.Checker = CheckerUnit
	}
	flags.Var(&info.Kind, "kind", "[inspect,ssa,codegen]")
	if info.Kind == "" {
		info.Kind = KindInspect
	}
	flags.BoolVar(&info.Cmd, "cmd", true, "create main file")
	flags.BoolVar(&info.Plugin, "plugin", false, "create golangci-lint plugin")
	flags.StringVar(&info.Pkg, "pkg", "", "package name")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	return flags.Args(), nil
}

func (s *Skeleton) run(info *Info) error {
	fsys, err := new(Generator).Run(info)
	if err != nil {
		return err
	}

	prompt := &Prompt{
		Output:    s.Output,
		ErrOutput: s.ErrOutput,
		Input:     s.Input,
	}

	dst := filepath.Join(s.Dir, info.Pkg)
	if err := CreateDir(prompt, dst, fsys); err != nil {
		return err
	}
	return nil
}
