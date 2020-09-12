package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"go/format"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/tools/txtar"
)

//go:generate go run tools/txtar/main.go -strip "_template/" _template template.go

func main() {
	var s Skeleton
	flag.BoolVar(&s.OverWrite, "overwrite", false, "overwrite all file")
	flag.BoolVar(&s.Cmd, "cmd", true, "create cmd directory")
	flag.StringVar(&s.Checker, "checker", "unit", "checker which is used in main.go (unit,single,multi)")
	flag.BoolVar(&s.Plugin, "plugin", true, "create plugin directory")
	flag.StringVar(&s.Type, "type", "inspect", "type of skeleton code (inspect|ssa)")
	flag.StringVar(&s.ImportPath, "path", "", "import path")
	flag.Parse()
	s.ExeName = os.Args[0]
	s.Args = flag.Args()

	switch s.Checker {
	case "unit", "single", "multi":
		// noop
	default:
		s.Checker = "unit"
	}

	if err := s.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

type Skeleton struct {
	ExeName    string
	Args       []string
	Dir        string
	ImportPath string
	OverWrite  bool
	Cmd        bool
	Checker    string
	Plugin     bool
	Mode       Mode
	Type       string
}

type Mode int

const (
	ModeRemoveAndCreateNew Mode = iota
	ModeConfirm
	ModeCreateNewFile
)

type TemplateData struct {
	Pkg        string
	ImportPath string
	Cmd        bool
	Plugin     bool
	Checker    string
	Type       string
}

func (s *Skeleton) Run() error {

	td := &TemplateData{
		Cmd:     s.Cmd,
		Plugin:  s.Plugin,
		Checker: s.Checker,
		Type:    s.Type,
	}

	if len(s.Args) < 1 {
		if s.ImportPath != "" {
			s.Dir = s.ImportPath
			td.Pkg = path.Base(s.ImportPath)
		} else {
			return errors.New("package must be specified")
		}
	} else {
		s.Dir = s.Args[0]
		td.Pkg = path.Base(s.Args[0])
	}

	switch s.Type {
	case "inspect", "ssa":
	default:
		return fmt.Errorf("unexpected type: %s")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	td.ImportPath = s.importPath(cwd)

	if td.ImportPath == "" {
		const format = "%s must be executed in GOPATH or -path option must be specified"
		return fmt.Errorf(format, s.ExeName)
	}

	exist, err := isExist(s.Dir)
	if err != nil {
		return err
	}
	if exist && !s.OverWrite {
		if exit := s.selectMode(s.Dir); exit {
			return nil
		}
	}

	if err := s.createAll(td); err != nil {
		return err
	}

	return nil
}

func (s *Skeleton) importPath(cwd string) string {

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
			return path.Join(filepath.ToSlash(rel), filepath.ToSlash(s.Dir))
		}
	}

	return ""
}

func (s *Skeleton) selectMode(dir string) bool {
	fmt.Printf("%s already exist, remove?\n", dir)
	fmt.Println("[1] No(Exit)")
	fmt.Println("[2] Remove and create new directory")
	fmt.Println("[3] Overwrite existing files with confirmation")
	fmt.Println("[4] Create new files only")
	fmt.Print("(default is 1) >")
	var m string
	fmt.Scanln(&m)
	switch m {
	case "2":
		s.Mode = ModeRemoveAndCreateNew
	case "3":
		s.Mode = ModeConfirm
	case "4":
		s.Mode = ModeCreateNewFile
	default:
		// exit
		return true
	}
	return false
}

func (s *Skeleton) createAll(td *TemplateData) error {

	if s.Mode == ModeRemoveAndCreateNew {
		if err := os.RemoveAll(s.Dir); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, td); err != nil {
		return err
	}

	ar := txtar.Parse(buf.Bytes())
	for _, f := range ar.Files {
		if err := s.createFile(f); err != nil {
			return err
		}
	}

	return nil
}

func (s *Skeleton) createFile(f txtar.File) (rerr error) {
	if len(bytes.TrimSpace(f.Data)) == 0 {
		return nil
	}

	path := filepath.Join(s.Dir, filepath.FromSlash(f.Name))

	exist, err := isExist(path)
	if err != nil {
		return err
	}

	if exist {
		switch s.Mode {
		case ModeConfirm:
			fmt.Printf("%s already exit, replace? [y/N] >", path)
			var yn string
			fmt.Scanln(&yn)
			switch strings.ToLower(yn) {
			case "y", "yes":
				// continue
			default:
				// skip
				fmt.Println("skip", path)
				return nil
			}
		case ModeCreateNewFile:
			// skip
			fmt.Println("skip", path)
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := w.Close(); err != nil && rerr == nil {
			rerr = err
		}
	}()

	// format a go file
	data := f.Data
	if filepath.Ext(path) == ".go" {
		data, err = format.Source(data)
		if err != nil {
			return err
		}
	}

	r := bytes.NewReader(data)
	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	fmt.Println("create", path)

	return nil
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
