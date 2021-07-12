package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
)

//go:embed version.txt
var version string

func main() {
	os.Exit(skeleton.Main(strings.TrimSpace(version), os.Args[1:]))
}
