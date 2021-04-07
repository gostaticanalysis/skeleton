package cmd

type Kind int

const (
	KindInspect Kind = iota
	KindSSA
	KindCodegen
)

func ParseKind(s string) Kind {
	switch s {
	case "ssa":
		return KindSSA
	case "codegen":
		return KindCodegen
	default:
		return KindInspect
	}
}
