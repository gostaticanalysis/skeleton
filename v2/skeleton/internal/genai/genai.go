package genai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Generator struct {
	Model string
}

func New() *Generator {
	return &Generator{
		Model: openai.GPT4o,
	}
}

func (g *Generator) Do(ctx context.Context, client *openai.Client, w io.Writer, inst *Instruction) error {
	var prompt bytes.Buffer
	if err := WritePrompt(&prompt, inst); err != nil {
		return fmt.Errorf("genai.Generate: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: g.Model,
		Messages: []openai.ChatCompletionMessage{{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt.String(),
		}},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("genai.Generate: %w", err)
	}

	if len(resp.Choices) == 0 {
		return errors.New("genai.Generate: cannot genearte code")
	}

	src := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(resp.Choices[0].Message.Content, "```go"), "```"))
	if _, err := fmt.Fprint(w, src); err != nil {
		return fmt.Errorf("genai.Generator.Do: %w", err)
	}

	return nil

}

func Generate(ctx context.Context, client *openai.Client, w io.Writer, inst *Instruction) error {
	return New().Do(ctx, client, w, inst)
}
