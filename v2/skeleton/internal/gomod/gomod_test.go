package gomod_test

import (
	"fmt"
	"testing"

	"github.com/gostaticanalysis/skeleton/v2/skeleton/internal/gomod"
	"github.com/tenntenn/golden"
)

func TestParentModule(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path    string
		wantErr bool
	}{
		"exisit":   {"example.com/example", false},
		"noexisit": {"", true},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			if tt.path != "" {
				modfile := fmt.Sprintf("-- go.mod --\nmodule %s\ngo 1.18", tt.path)
				golden.DirInit(t, dir, modfile)
			}
			_, got, err := gomod.ParentModule(dir)

			switch {
			case !tt.wantErr && err != nil:
				t.Error("unexpected error:", err)
			case tt.wantErr && err == nil:
				t.Error("expected error did not occur")
			}

			if err == nil && got != tt.path {
				t.Errorf("want %s but got %s", tt.path, got)
			}
		})
	}
}
