package skeleton

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// DefaultPrompt is default Prompt.
var DefaultPrompt = &Prompt{
	Input:     os.Stdin,
	Output:    os.Stdout,
	ErrOutput: os.Stderr,
}

// Prompt receive input from a user.
type Prompt struct {
	Input     io.Reader
	Output    io.Writer
	ErrOutput io.Writer
}

// Choose shows options and user would select a option from these options.
// Choose returns index of the selected option.
func (p *Prompt) Choose(description string, opts []string, prompt string) (int, error) {
	if _, err := fmt.Fprintln(p.Output, description); err != nil {
		return 0, err
	}

	optfmt := fmt.Sprintf("[%%%dd] %%s\n", len(strconv.Itoa(len(opts))))
	for i := range opts {
		if _, err := fmt.Fprintf(p.Output, optfmt, i+1, opts[i]); err != nil {
			return 0, err
		}
	}
	for {
		if _, err := fmt.Fprint(p.Output, prompt); err != nil {
			return 0, err
		}

		var s string
		if _, err := fmt.Fscanln(p.Input, &s); err != nil {
			return 0, err
		}

		n, err := strconv.Atoi(s)
		if err == nil && n >= 1 && n <= len(opts) {
			return n - 1, nil
		}

		if _, err := fmt.Fprintf(p.ErrOutput, "%s is invalid option\n", s); err != nil {
			return 0, err
		}
	}
}

func (p *Prompt) YesNo(description string, defaultVal bool, prompt rune) (bool, error) {
	if _, err := fmt.Fprintln(p.Output, description); err != nil {
		return false, err
	}

	for {
		ynfmt := "[y/N]%c"
		if defaultVal {
			ynfmt = "[Y/n]%c"
		}
		if _, err := fmt.Fprintf(p.Output, ynfmt, prompt); err != nil {
			return false, err
		}

		var s string
		if _, err := fmt.Fscanln(p.Input, &s); err != nil {
			return false, err
		}
		switch strings.ToUpper(s) {
		case "Y", "YES", "OK":
			return true, nil
		case "N", "no":
			return false, nil
		case "":
			return defaultVal, nil
		}

		if _, err := fmt.Fprintf(p.ErrOutput, "%s is invalid option", s); err != nil {
			return false, err
		}
	}
}
