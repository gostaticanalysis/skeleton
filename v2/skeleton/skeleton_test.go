package skeleton_test

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/tenntenn/golden"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
	os.Setenv("SKELETON_PREFIX", "")
}

func TestSkeletonRun(t *testing.T) {
	t.Parallel()
	F := golden.TxtarWith
	cases := map[string]struct {
		dirinit string
		args    string
		input   string

		wantExitCode int
		wantOutput   string
		wantGoTest   bool
	}{
		"nooption":              {"", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"overwrite-cancel":      {F(t, "example/go.mod", "// empty"), "example.com/example", "1\n", skeleton.ExitSuccess, "", false},
		"overwrite-force":       {F(t, "example/go.mod", "// empty"), "example.com/example", "2\n", skeleton.ExitSuccess, "", true},
		"overwrite-confirm-yes": {F(t, "example/go.mod", "// empty"), "example.com/example", "3\ny\n", skeleton.ExitSuccess, "", true},
		"overwrite-confirm-no":  {F(t, "example/go.mod", "// empty"), "example.com/example", "3\nn\n", skeleton.ExitSuccess, "", false},
		"overwrite-newonly":     {F(t, "example/go.mod", "// empty"), "example.com/example", "4\n", skeleton.ExitSuccess, "", false},
		"plugin":                {"", "-plugin example.com/example", "", skeleton.ExitSuccess, "", true},
		"nocmd":                 {"", "-cmd=false example.com/example", "", skeleton.ExitSuccess, "", true},
		"onlypkgname":           {"", "example", "", skeleton.ExitSuccess, "", true},
		"version":               {"", "-v", "", skeleton.ExitSuccess, "skeleton version\n", false},
		"kind-inspect":          {"", "-kind inspect example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-ssa":              {"", "-kind ssa example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-codegen":          {"", "-kind codegen example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-packages":         {"", "-kind packages example.com/example", "", skeleton.ExitSuccess, "", true},
	}

	if flagUpdate {
		golden.RemoveAll(t, "testdata")
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			if tt.dirinit != "" {
				golden.DirInit(t, dir, tt.dirinit)
			}

			var out, errout bytes.Buffer
			s := &skeleton.Skeleton{
				Dir:       dir,
				Output:    &out,
				ErrOutput: &errout,
				Input:     strings.NewReader(tt.input),
			}

			args := strings.Split(tt.args, " ")
			gotExitCode := s.Run(name, args)

			if gotExitCode != tt.wantExitCode {
				t.Errorf("exit code want %d got %d", tt.wantExitCode, gotExitCode)
			}

			if tt.wantExitCode == 0 && errout.String() != "" {
				t.Error("exit code want 0 but error messages are outputed", errout.String())
			}

			if tt.wantOutput != "" && out.String() != tt.wantOutput {
				t.Errorf("output want %s got %s", tt.wantOutput, out.String())
			}

			got := golden.Txtar(t, s.Dir)

			if tt.wantGoTest {
				entries, err := os.ReadDir(dir)
				if err != nil {
					t.Fatal("unexpected error:", err)
				}

				if len(entries) > 0 {
					skeletondir := filepath.Join(dir, entries[0].Name())
					gomodtidy(t, skeletondir)
					gotest(t, name, skeletondir)
				}
			}

			if flagUpdate {
				golden.Update(t, "testdata", name, got)
				return
			}

			if diff := golden.Diff(t, "testdata", name, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func gomodtidy(t *testing.T, dir string) {
	t.Helper()
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("go mod tidy: unexpected error: %s with:\n%s", err, &stderr)
	}
}

var (
	timeRegexp = regexp.MustCompile(`([\(\t])([0-9.]+s)(\)?)`)
	hexRegexp  = regexp.MustCompile(`(\()(0x[0-9a-f]+)(\))`)
)

func gotest(t *testing.T, name, dir string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "test")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	t.Log("exec", cmd)

	if err := cmd.Run(); err != nil && !errors.As(err, new(*exec.ExitError)) {
		t.Fatal("unexpected error:", err)
	}

	got := stdout.String() + stderr.String()
	got = timeRegexp.ReplaceAllString(got, "${1}0000s${3}")
	got = hexRegexp.ReplaceAllString(got, "${1}0x0000${3}")

	goldenname := name + "-go-test"
	if flagUpdate {
		golden.Update(t, "testdata", goldenname, got)
		return
	}

	if diff := golden.Diff(t, "testdata", goldenname, got); diff != "" {
		t.Error(diff)
	}
}
