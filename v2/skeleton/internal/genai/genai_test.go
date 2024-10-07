package genai_test

import (
	"bytes"
	"context"
	"fmt"
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
		//Details: []string{},
		Tests: `package a

import "log"

// The tool can find function callings of log.Fatal.
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

	fmt.Print(buf.String())
}
