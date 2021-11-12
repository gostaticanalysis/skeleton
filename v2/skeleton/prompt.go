package skeleton

import (
	"os"

	"github.com/gostaticanalysis/skeletonkit"
)

// DefaultPrompt is default Prompt.
var DefaultPrompt = &Prompt{
	Input:     os.Stdin,
	Output:    os.Stdout,
	ErrOutput: os.Stderr,
}

// Prompt receive input from a user.
type Prompt = skeletonkit.Prompt
