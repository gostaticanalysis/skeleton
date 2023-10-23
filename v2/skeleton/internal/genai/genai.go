package genai

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

func Generate(ctx context.Context, client *openai.Client, w io.Writer, inst *Instruction) error {
	var prompt bytes.Buffer
	if err := WritePrompt(&prompt, inst); err != nil {
		return fmt.Errorf("genai.Generate: %w", err)
	}

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
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

	if _, err := fmt.Fprint(w, resp.Choices[0].Message.Content); err != nil {
		return fmt.Errorf("genai.Generate: %w", err)
	}

	return nil
}
