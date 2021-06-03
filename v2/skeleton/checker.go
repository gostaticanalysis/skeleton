package skeleton

import "flag"

type Checker string

const (
	CheckerUnit   Checker = "unit"
	CheckerSingle Checker = "single"
	CheckerMulti  Checker = "multi"
)

var _ flag.Value = (*Checker)(nil)

// String returns "single", "multi" or "unit".
func (ch Checker) String() string {
	switch ch {
	case CheckerSingle:
		return "single"
	case CheckerMulti:
		return "multi"
	default:
		return "unit"
	}
}

func (ch *Checker) Set(s string) error {
	switch s {
	case "single":
		*ch = CheckerSingle
	case "multi":
		*ch = CheckerMulti
	default:
		*ch = CheckerUnit
	}
	return nil
}
