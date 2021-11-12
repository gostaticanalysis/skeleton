package skeleton

import (
	"embed"
	"text/template"

	"github.com/gostaticanalysis/skeletonkit"
)

//go:embed _template/*
var tmplFS embed.FS

// DefaultTemplate is default template for skeleton.
var DefaultTemplate *template.Template

// DefaultFuncMap is default FuncMap for a template.
// Deprecated: should use skeletonkit.TemplateWithFuncs
var DefaultFuncMap = skeletonkit.DefaultFuncMap

func init() {
	DefaultTemplate = template.Must(skeletonkit.ParseTemplate(tmplFS, "skeleton", "_template"))
}
