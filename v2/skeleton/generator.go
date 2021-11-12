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
	return skeletonkit.ExecuteTemplate(g.template(), info)
}

func (g *Generator) template() *template.Template {
	if g.Template != nil {
		return g.Template
	}
	return DefaultTemplate
}
