package skeleton

import (
	"bytes"
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/imports"
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
	ar := txtar.Parse(buf.Bytes())
	for i := range ar.Files {
		if filepath.Ext(ar.Files[i].Name) != ".go" || len(ar.Files[i].Data) == 0 {
			continue
		}
		opt := &imports.Options{
			FormatOnly: true,
		}
		src, err := imports.Process(ar.Files[i].Name, ar.Files[i].Data, opt)
		if err != nil {
			return nil, err
		}
		ar.Files[i].Data = src
	}
	return txtarfs.As(ar), nil
}

func (g *Generator) template() *template.Template {
	if g.Template != nil {
		return g.Template
	}
	return DefaultTemplate
}
