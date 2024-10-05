package genai_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/gostaticanalysis/skeleton/v2/skeleton/internal/genai"
	openai "github.com/sashabaranov/go-openai"
)

func TestGenerate(t *testing.T) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()
	var buf bytes.Buffer

	inst := &genai.Instruction{
		Pkg: "logfatal",
		Details: []string{
			"The tool can find function callings of log.Fatal.",
		},
		Tests: `package a

import "log"

func f() {
	log.Fatal("error") // want "NG"
	fatal := log.Fatal
	fatal("error") // want "NG"
	println() // OK
}
`,
	}

	if err := genai.Generate(ctx, client, &buf, inst); err != nil {
		t.Fatal("unexpected error:", err)
	}

	//t.Log(buf.String())
}
