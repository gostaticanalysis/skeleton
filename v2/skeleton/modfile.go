package skeleton

import (
	"errors"
	"fmt"
	"go/build"
	"strings"

	"golang.org/x/mod/modfile"
)

func modinit(path string) (string, error) {
	var mf modfile.File
	if err := mf.AddModuleStmt(path); err != nil {
		return "", fmt.Errorf("create go.mod: %w", err)
	}

	gov, err := goVersion()
	if err != nil {
		return "", fmt.Errorf("create go.mod: %w", err)
	}

	if err := mf.AddGoStmt(gov); err != nil {
		return "", fmt.Errorf("create go.mod: %w", err)
	}

	b, err := mf.Format()
	if err != nil {
		return "", fmt.Errorf("create go.mod: %w", err)
	}

	return string(b), nil
}

func goVersion() (string, error) {
	tags := build.Default.ReleaseTags
	for i := len(tags) - 1; i >= 0; i-- {
		version := tags[i]
		if strings.HasPrefix(version, "go") && modfile.GoVersionRE.MatchString(version[2:]) {
			return version[2:], nil
		}
	}
	return "", errors.New("there are not valid go version")
}
