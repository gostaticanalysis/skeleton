package skeleton

import (
	"io/fs"
	"text/template"

	"github.com/gostaticanalysis/skeletonkit"
)

type Generator struct {
	Template *template.Template
}

func (g *Generator) Run(info *Info) (fs.FS, error) {
	tmpl, err := g.template(info)
	if err != nil {
		return nil, err
	}
	return skeletonkit.ExecuteTemplate(tmpl, info)
}

func (g *Generator) template(info *Info) (*template.Template, error) {
	if g.Template != nil {
		return g.Template, nil
	}
	return parseTemplate(info.Kind)
}
