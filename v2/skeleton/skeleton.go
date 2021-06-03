package skeleton

import (
	"bytes"
	"context"
	"io"
	"io/fs"
	"os"
	"text/template"

	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

type Generator struct {
	Stdout, Stderr io.Writer
	Template       *template.Template
}

func (g *Generator) Run(ctx context.Context, info *Info) (fs.FS, error) {
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
