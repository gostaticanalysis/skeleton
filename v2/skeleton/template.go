package skeleton

import (
	"embed"
	"path"
	"text/template"

	"github.com/gostaticanalysis/skeletonkit"
	"golang.org/x/mod/semver"
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

func parseTemplate(info *Info) (*template.Template, error) {
	dir := info.Kind.String()
	if dir != "packages" && go118(info.GoVersion) {
		dir += "_go118"
	}
	return skeletonkit.ParseTemplate(tmplFS, "skeleton", path.Join("_template", dir))
}

func go118(v string) bool {
	return v != "" && semver.Compare("v"+v, "v1.18") >= 0
}
