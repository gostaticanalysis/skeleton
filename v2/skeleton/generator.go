package skeleton

import (
	"bytes"
	"io/fs"
	"text/template"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

type Generator struct {
	Template *template.Template
}

func (g *Generator) Run(info *Info) (fs.FS, error) {
	var buf bytes.Buffer
	if err := g.template().Execute(&buf, info); err != nil {
		return nil, err
	}
	return txtarfs.As(txtar.Parse(buf.Bytes())), nil
}

func (g *Generator) template() *template.Template {
	if g.Template != nil {
		return g.Template
	}
	return DefaultTemplate
}
