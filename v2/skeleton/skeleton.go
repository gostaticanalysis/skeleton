package skeleton

import (
	"io"
	"os"
)

type Generator struct {	
	Stdout, Stderr io.Writer
	Template       *text.Template
}

func (g *Generator) template() *text.Template {
	if g.Template != nil {
		return g.Template
	}
	return DefaultTemplate
}

func (g *Generator) stdout() io.Writer {
	if g.Stdout != nil {
		return g.Stdout
	}
	return os.Stdout
}

func (g *Generator) stderr() io.Writer {
	if g.Stderr != nil {
		return g.Stderr
	}
	return os.Stderr
}

func (g *Generator) Run(ctx context.Context, info *Info) (fs.Fs, error) {
	var buf bytes.Buffer
}
