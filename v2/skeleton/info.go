package skeleton

type Info struct {
	Kind        Kind
	Checker     Checker
	Pkg         string
	Path        string
	Cmd         bool
	OmmitPlugin bool
}
