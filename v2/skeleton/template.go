package skeleton

import (
	"embed"
	"text/template"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

//go:embed _template/*
var tmplFS embed.FS

// DefaultTemplate is
var DefaultTemplate *text.Template

func init() {
	ar, err := txtarfs.From(tmplFS)
	if err != nil {
		panic(err)
	}
	strTmpl := string(txtar.Format(ar))
	DefaultTemplate = template.Must(template.New("skeleton").Delims("@@", "@@").Parse(strTmpl))
}
