package skeleton

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/gostaticanalysis/skeleton/v2/skeleton/internal/gomod"
	"github.com/gostaticanalysis/skeletonkit"
	"golang.org/x/mod/modfile"
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
	flags, err := s.parseFlag(args, &info)
	if err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return ExitError
	}

	info.Path = flags.Arg(0)
	if !info.GoMod {
		parentModule, err := gomod.ParentModule(s.Dir)
		if err != nil {
			fmt.Fprintln(s.ErrOutput, "Error:", err)
			return ExitError
		}
		info.Path = path.Join(parentModule, info.Path)
	} else if prefix := os.Getenv("SKELETON_PREFIX"); prefix != "" {
		info.Path = path.Join(prefix, info.Path)
	}

	// allow package name only
	if module.CheckImportPath(info.Path) != nil {
		flags.Usage()
		return ExitError
	}

	if info.Pkg == "" {
		info.Pkg = path.Base(info.Path)
	}

	if err := s.run(&info); err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return ExitError
	}

	return ExitSuccess
}

func (s *Skeleton) parseFlag(args []string, info *Info) (*flag.FlagSet, error) {
	flags := flag.NewFlagSet("skeleton", flag.ContinueOnError)
	flags.SetOutput(s.ErrOutput)
	flags.Usage = func() {
		fmt.Fprintln(s.ErrOutput, "skeleton [-checker,-kind,-cmd,-plugin] example.com/path")
		flags.PrintDefaults()
	}
	flags.Var(&info.Checker, "checker", "[unit,single,multi]")

	flags.Var(&info.Kind, "kind", "[inspect,ssa,codegen,packages]")

	flags.BoolVar(&info.Cmd, "cmd", true, "create main file")
	flags.BoolVar(&info.Plugin, "plugin", false, "create golangci-lint plugin")
	flags.StringVar(&info.Pkg, "pkg", "", "package name")
	flags.BoolVar(&info.GoMod, "gomod", true, "create a go.mod file")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if info.Kind == "" {
		info.Kind = KindInspect
	}

	if info.Checker == "" {
		switch info.Kind {
		case KindCodegen:
			info.Checker = CheckerSingle
		default:
			info.Checker = CheckerUnit
		}
	}

	return flags, nil
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
	opts := []skeletonkit.CreatorOption{
		skeletonkit.CreatorWithEmpty(true),
		skeletonkit.CreatorWithSkipFunc(func(p string, d fs.DirEntry) bool {
			switch {
			case !info.Plugin && path.Base(p) == "plugin":
				return true
			case !info.GoMod && path.Base(p) == "go.mod":
				return true
			}
			return false // no skip
		}),
	}
	if err := skeletonkit.CreateDir(prompt, dst, fsys, opts...); err != nil {
		return err
	}
	return nil
}

func ParentModule(dir string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("can not get the parent module: %w", err)
	}

	gomodfile := strings.TrimSpace(stdout.String())
	moddata, err := os.ReadFile(gomodfile)
	if err != nil {
		return "", fmt.Errorf("cat not read the go.mod of the parent module: %w", err)
	}

	gomod, err := modfile.Parse(gomodfile, moddata, nil)
	if err != nil {
		return "", fmt.Errorf("cat parse the go.mod of the parent module: %w", err)
	}

	return gomod.Module.Mod.Path, nil
}
