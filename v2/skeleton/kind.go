package skeleton

import "flag"

// Kind represents kind of skeleton codes.
// Kind implements flag.Value.
type Kind string

var _ flag.Value = (*Kind)(nil)

const (
	KindInspect  Kind = "inspect"
	KindSSA      Kind = "ssa"
	KindCodegen  Kind = "codegen"
	KindPackages Kind = "packages"
)

func (k Kind) String() string {
	switch k {
	case KindSSA:
		return "ssa"
	case KindCodegen:
		return "codegen"
	case KindPackages:
		return "packages"
	default:
		return "inspect"
	}
}

// "ssa" -> KindSSA, "codegen" -> KindCodegen, "packages" -> KindPackages otherwise KindInspect.
func (k *Kind) Set(s string) error {
	switch s {
	case "ssa":
		*k = KindSSA
	case "codegen":
		*k = KindCodegen
	case "packages":
		*k = KindPackages
	default:
		*k = KindInspect
	}
	return nil
}
