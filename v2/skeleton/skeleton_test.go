package skeleton_test

import (
	"bytes"
	"flag"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/josharian/txtarfs"
	"golang.org/x/tools/txtar"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func TestSkeletonRun(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		args  string
		input string

		wantExitCode int
	}{
		"nooption": {"example.com/example", "", skeleton.ExitSuccess},
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var out, errout bytes.Buffer

			s := &skeleton.Skeleton{
				Dir:       t.TempDir(),
				Output:    &out,
				ErrOutput: &errout,
				Input:     strings.NewReader(tt.input),
			}

			args := strings.Split(tt.args, " ")
			gotExitCode := s.Run(name, args)

			if gotExitCode != tt.wantExitCode {
				t.Errorf("exit code want %d got %d", tt.wantExitCode, gotExitCode)
			}

			got := txtarDir(t, s.Dir)

			if flagUpdate {
				writeTestData(t, name, got)
				return
			}

			want := readTestData(t, name)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func txtarDir(t *testing.T, dir string) string {
	t.Helper()
	ar, err := txtarfs.From(os.DirFS(dir))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	return string(txtar.Format(ar))
}

func writeTestData(t *testing.T, name, v string) {
	t.Helper()
	f, err := os.Create(filepath.Join("testdata", name + ".golden"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Fatal("unexpected error", err)
		}
	}()

	if _, err := io.Copy(f, strings.NewReader(v)); err != nil {
		t.Fatal("unexpected error", err)
	}
}

func readTestData(t *testing.T, name string) string {
	t.Helper()
	got, err := os.ReadFile(filepath.Join("testdata", name + ".golden"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	return string(got)
}
