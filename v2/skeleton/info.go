package skeleton

type Info struct {
	Kind       Kind
	Checker    Checker
	Pkg        string
	ImportPath string
	Cmd        bool
	Plugin     bool
}
