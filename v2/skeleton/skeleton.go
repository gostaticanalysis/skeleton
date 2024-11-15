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

	"github.com/gostaticanalysis/skeletonkit"
	"golang.org/x/mod/module"

	"github.com/gostaticanalysis/skeleton/v2/skeleton/internal/gomod"
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
	GoVersion string
}

func Main(version string, args []string) int {
	s := &Skeleton{
		Dir:       ".",
		Output:    os.Stdout,
		ErrOutput: os.Stderr,
		Input:     os.Stdin,
	}

	gover, err := goVersion(s.Dir)
	if err != nil {
		fmt.Fprintln(s.ErrOutput, "Error:", err)
		return ExitError
	}

	if strings.HasPrefix(gover, "devel ") {
		// The devel format is like following.
		// devel go1.24-xxxxxxxxxx Day Date Mon Time Arch"
		gover = strings.Split(strings.TrimPrefix(gover, "devel "), "-")[0]
	}

	s.GoVersion = gover

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

	info.GoVersion = s.GoVersion

	info.Path = flags.Arg(0)
	if !info.GoMod {
		importpath, err := s.withoutGoMod(info.Path)
		if err != nil {
			fmt.Fprintln(s.ErrOutput, "Error:", err)
			return ExitError
		}
		info.Path = importpath
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
			case !info.GoMod && isGoMod(p):
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

func (s *Skeleton) withoutGoMod(p string) (string, error) {
	moddir, modpath, err := gomod.ParentModule(s.Dir)
	if err != nil {
		return "", err
	}

	wd, err := filepath.EvalSymlinks(s.Dir)
	if err != nil {
		return "", err
	}

	wd, err = filepath.Abs(wd)
	if err != nil {
		return "", err
	}

	moddir, err = filepath.EvalSymlinks(moddir)
	if err != nil {
		return "", err
	}

	moddir, err = filepath.Abs(moddir)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(moddir, wd)
	if err != nil {
		return "", err
	}

	return path.Join(modpath, filepath.ToSlash(rel), p), nil
}

func isGoMod(p string) bool {
	return path.Base(p) == "go.mod" &&
		!strings.Contains(p, "testdata/")
}

func goVersion(dir string) (string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("go", "env", "GOVERSION")
	cmd.Dir = dir
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
