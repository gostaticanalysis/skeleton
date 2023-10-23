package genai

import (
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

var (
	//go:embed prompt_template.txt
	promptTmplFile string
	promptTmpl     = template.Must(template.New("prompt").Parse(promptTmplFile))
)

type Instruction struct {
	Pkg     string
	Details []string
}

func WritePrompt(w io.Writer, inst *Instruction) error {
	if err := promptTmpl.Execute(w, inst); err != nil {
		return fmt.Errorf("create prompt: %w", err)
	}
	return nil
}
