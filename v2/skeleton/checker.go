package skeleton

import "fmt"

type Checker int

const (
	CheckerUnit Checker = iota
	CheckerSingle
	CheckerMulti
)

var _ fmt.Stringer = CheckerUnit

// String implements fmt.Stringer.
// It will return "single", "multi" or "unit".
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
