package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func main() {
	var s Skeleton
	flag.BoolVar(&s.Cmd, "cmd", true, "create cmd directory")
	flag.BoolVar(&s.Plugin, "plugin", true, "create plugin directory")
	flag.StringVar(&s.ImportPath, "path", "", "import path")
	flag.Parse()
	s.ExeName = os.Args[0]
	s.Args = flag.Args()

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type PkgInfo struct {
	Pkg        string
	ImportPath string
}

type Skeleton struct {
	ExeName    string
	Args       []string
	Cmd        bool
	Plugin bool
	ImportPath string
}

func (s *Skeleton) Run() error {

	var info PkgInfo

	if len(s.Args) < 1 {
		if s.ImportPath != "" {
			info.Pkg = path.Base(s.ImportPath)
		} else {
			return errors.New("package must be specified")
		}
	} else {
		info.Pkg = s.Args[0]
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	info.ImportPath = s.importPath(cwd, &info)

	if info.ImportPath == "" {
		const format = "%s must be executed in GOPATH or -path option must be specified"
		return errors.Errorf(format, s.ExeName)
	}

	dir := filepath.Join(cwd, info.Pkg)
	if err := os.Mkdir(dir, 0777); err != nil {
		return err
	}

	src, err := os.Create(filepath.Join(dir, info.Pkg+".go"))
	if err != nil {
		return err
	}
	defer src.Close()
	if err := srcTempl.Execute(src, info); err != nil {
		return err
	}

	test, err := os.Create(filepath.Join(dir, info.Pkg+"_test.go"))
	if err != nil {
		return err
	}
	defer test.Close()
	if err := testTempl.Execute(test, info); err != nil {
		return err
	}

	testdata := filepath.Join(dir, "testdata", "src", "a")
	if err := os.MkdirAll(testdata, 0777); err != nil {
		return err
	}

	adotgo, err := os.Create(filepath.Join(testdata, "a.go"))
	if err != nil {
		return err
	}
	defer adotgo.Close()
	if err := adotgoTempl.Execute(adotgo, info); err != nil {
		return err
	}

	if s.Cmd {
		if err := s.createCmd(dir, &info); err != nil {
			return err
		}
	}

	if s.Plugin {
		if err := s.createPlugin(dir, &info); err != nil {
			return err
		}
	}

	return nil
}

func (s *Skeleton) importPath(cwd string, info *PkgInfo) string {

	if s.ImportPath != "" {
		return s.ImportPath
	}

	for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
		if gopath == "" {
			continue
		}

		src := filepath.Join(gopath, "src")
		if strings.HasPrefix(cwd, src) {
			rel, err := filepath.Rel(src, cwd)
			if err != nil {
				return ""
			}
			return path.Join(filepath.ToSlash(rel), info.Pkg)
		}
	}

	return ""
}

func (s *Skeleton) createCmd(dir string, info *PkgInfo) error {
	cmdDir := filepath.Join(dir, "cmd", info.Pkg)
	if err := os.MkdirAll(cmdDir, 0777); err != nil {
		return err
	}

	cmdMain, err := os.Create(filepath.Join(cmdDir, "main.go"))
	if err != nil {
		return err
	}
	defer cmdMain.Close()

	if err := cmdMainTempl.Execute(cmdMain, info); err != nil {
		return err
	}

	return nil
}

func (s *Skeleton) createPlugin(dir string, info *PkgInfo) error {
	pluginDir := filepath.Join(dir, "plugin", info.Pkg)
	if err := os.MkdirAll(pluginDir, 0777); err != nil {
		return err
	}

	pluginMain, err := os.Create(filepath.Join(pluginDir, "main.go"))
	if err != nil {
		return err
	}
	defer pluginMain.Close()

	if err := pluginMainTempl.Execute(pluginMain, info); err != nil {
		return err
	}

	return nil
}
