package skeleton_test

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/tenntenn/golden"

	"github.com/gostaticanalysis/skeleton/v2/skeleton"
	"github.com/gostaticanalysis/skeleton/v2/skeleton/internal/gomod"
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
	const noflags = ""
	cases := map[string]struct {
		goVersion string
		dir       string
		dirinit   string
		flags     string
		path      string
		input     string

		wantExitCode int
		wantOutput   string
		wantGoTest   bool
	}{
		"nooption":              {"", "", "", "", "example.com/example", noflags, skeleton.ExitSuccess, "", true},
		"overwrite-cancel":      {"", "", F(t, "example/go.mod", "// empty"), noflags, "example.com/example", "1\n", skeleton.ExitSuccess, "", false},
		"overwrite-force":       {"", "", F(t, "example/go.mod", "// empty"), noflags, "example.com/example", "2\n", skeleton.ExitSuccess, "", true},
		"overwrite-confirm-yes": {"", "", F(t, "example/go.mod", "// empty"), noflags, "example.com/example", "3\ny\n", skeleton.ExitSuccess, "", true},
		"overwrite-confirm-no":  {"", "", F(t, "example/go.mod", "// empty"), noflags, "example.com/example", "3\nn\n", skeleton.ExitSuccess, "", false},
		"overwrite-newonly":     {"", "", F(t, "example/go.mod", "// empty"), noflags, "example.com/example", "4\n", skeleton.ExitSuccess, "", false},
		"plugin":                {"", "", "", "-plugin", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"nocmd":                 {"", "", "", "-cmd=false", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"onlypkgname":           {"", "", "", noflags, "example", "", skeleton.ExitSuccess, "", true},
		"version":               {"", "", "", "-v", "", "", skeleton.ExitSuccess, "skeleton version\n", false},
		"kind-inspect":          {"", "", "", "-kind inspect", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-ssa":              {"", "", "", "-kind ssa", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-codegen":          {"", "", "", "-kind codegen", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-packages":         {"", "", "", "-kind packages", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"parent-module":         {"", "", F(t, "go.mod", "module example.com/example"), "-gomod=false", "sub", "", skeleton.ExitSuccess, "", true},
		"parent-module-deep":    {"", "sub", F(t, "go.mod", "module example.com/example", "sub/sub.go", "package sub"), "-gomod=false", "subsub", "", skeleton.ExitSuccess, "", true},
		"kind-inspect-copy-parent-gomod-to-testdata": {"", "", "", "-kind inspect -copy-parent-gomod", "example.com/example", "", skeleton.ExitSuccess, "", true},
		"kind-ssa-copy-parent-gomod-to-testdata":     {"", "", "", "-kind ssa -copy-parent-gomod", "example.com/example", "", skeleton.ExitSuccess, "", true},
	}

	if flagUpdate {
		golden.RemoveAll(t, "testdata")
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tmpdir := t.TempDir()
			if tt.dirinit != "" {
				golden.DirInit(t, tmpdir, tt.dirinit)
			}

			var out, errout bytes.Buffer
			s := &skeleton.Skeleton{
				Dir:       filepath.Join(tmpdir, tt.dir),
				Output:    &out,
				ErrOutput: &errout,
				Input:     strings.NewReader(tt.input),
				GoVersion: "1.17", // do not change it even if your go version is over 1.18
			}

			if tt.goVersion != "" {
				s.GoVersion = tt.goVersion
			}

			var args []string
			if tt.flags != "" {
				args = strings.Split(tt.flags, " ")
			}
			if tt.path != "" {
				args = append(args, tt.path)
			}
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

			if tt.wantGoTest && tt.path != "" {
				skeletondir := filepath.Join(s.Dir, path.Base(tt.path))
				modroot := modroot(t, skeletondir)
				execCmd(t, modroot, "go", "mod", "tidy")
				gotest(t, name, skeletondir)
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

func modroot(t *testing.T, dir string) string {
	t.Helper()
	modfile, err := gomod.ModFile(dir)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	return filepath.Dir(modfile)
}

var (
	timeRegexp = regexp.MustCompile(`([\(\t])([0-9.]+s)(\)?)`)
)

func gotest(t *testing.T, name, dir string) {
	t.Helper()

	var stdout, stderr bytes.Buffer
	cmd := exec.Command("go", "test")
	cmd.Dir = dir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil && !errors.As(err, new(*exec.ExitError)) {
		t.Fatal("unexpected error:", err)
	}

	got := stdout.String() + stderr.String()
	got = timeRegexp.ReplaceAllString(got, "${1}0000s${3}")

	goldenname := name + "-go-test"
	if flagUpdate {
		golden.Update(t, "testdata", goldenname, got)
		return
	}

	if diff := golden.Diff(t, "testdata", goldenname, got); diff != "" {
		t.Error(diff)
	}
}

func execCmd(t *testing.T, dir, cmd string, args ...string) io.Reader {
	t.Helper()
	var stdout, stderr bytes.Buffer
	_cmd := exec.Command(cmd, args...)
	_cmd.Stdout = &stdout
	_cmd.Stderr = &stderr
	_cmd.Dir = dir
	t.Log("exec", cmd, strings.Join(args, " "))
	if err := _cmd.Run(); err != nil {
		t.Fatal(err, "\n", &stderr)
	}
	return &stdout
}
