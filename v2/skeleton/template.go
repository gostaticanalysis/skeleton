package skeleton

import (
	"embed"
	"path"
	"text/template"

	"github.com/gostaticanalysis/skeletonkit"
)

//go:embed _template/*
var tmplFS embed.FS

// DefaultTemplate is default template for skeleton.
// Deprecated: should use skeletonkit.
var DefaultTemplate *template.Template

// DefaultFuncMap is default FuncMap for a template.
// Deprecated: should use skeletonkit.TemplateWithFuncs
var DefaultFuncMap = skeletonkit.DefaultFuncMap

func init() {
	// for backward compatibility
	DefaultTemplate = template.Must(skeletonkit.ParseTemplate(tmplFS, "skeleton", "_template/inspect"))
}

func parseTemplate(kind Kind) (*template.Template, error) {
	return skeletonkit.ParseTemplate(tmplFS, "skeleton", path.Join("_template", kind.String()))
}
