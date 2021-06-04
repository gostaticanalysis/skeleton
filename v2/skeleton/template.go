package skeleton

import (
	"embed"
	"io/fs"
	"text/template"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

//go:embed _template/*
var tmplFS embed.FS

// DefaultTemplate is default template for skeleton.
var DefaultTemplate *template.Template

// DefaultFuncMap is default FuncMap for a template.
var DefaultFuncMap = template.FuncMap{
	"gomod": func() string {
		return "go.mod"
	},
	"gomodinit": func(path string) string {
		f, err := modinit(path)
		if err != nil {
			panic(err)
		}
		return f
	},
}

func init() {
	fsys, err := fs.Sub(tmplFS, "_template")
	if err != nil {
		panic(err)
	}
	ar, err := txtarfs.From(fsys)
	if err != nil {
		panic(err)
	}
	strTmpl := string(txtar.Format(ar))
	DefaultTemplate = template.Must(template.New("skeleton").Delims("@@", "@@").Funcs(DefaultFuncMap).Parse(strTmpl))
}
