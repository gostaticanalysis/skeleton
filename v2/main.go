package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"golang.org/x/mod/module"
)

//go:embed version.txt
var version string

func main() {
	os.Exit(skeleton.Main(version, os.Args[1:]))
}
