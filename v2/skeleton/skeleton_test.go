package skeleton_test

import (
	"bytes"
	"flag"
	"os"
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
	}{
		"nooption":              {"", "example.com/example", "", skeleton.ExitSuccess, ""},
		"overwrite-cancel":      {F(t, "example/go.mod", "// empty"), "example.com/example", "1\n", skeleton.ExitSuccess, ""},
		"overwrite-force":       {F(t, "example/go.mod", "// empty"), "example.com/example", "2\n", skeleton.ExitSuccess, ""},
		"overwrite-confirm-yes": {F(t, "example/go.mod", "// empty"), "example.com/example", "3\ny\n", skeleton.ExitSuccess, ""},
		"overwrite-confirm-no":  {F(t, "example/go.mod", "// empty"), "example.com/example", "3\nn\n", skeleton.ExitSuccess, ""},
		"overwrite-newonly":     {F(t, "example/go.mod", "// empty"), "example.com/example", "4\n", skeleton.ExitSuccess, ""},
		"plugin":                {"", "-plugin example.com/example", "", skeleton.ExitSuccess, ""},
		"nocmd":                 {"", "-cmd=false example.com/example", "", skeleton.ExitSuccess, ""},
		"onlypkgname":           {"", "example", "", skeleton.ExitSuccess, ""},
		"version":               {"", "-v", "", skeleton.ExitSuccess, "skeleton version\n"},
		"kind-inspect":          {"", "-kind inspect example.com/example", "", skeleton.ExitSuccess, ""},
		"kind-ssa":              {"", "-kind ssa example.com/example", "", skeleton.ExitSuccess, ""},
		"kind-codegen":          {"", "-kind codegen example.com/example", "", skeleton.ExitSuccess, ""},
		"kind-packages":         {"", "-kind packages example.com/example", "", skeleton.ExitSuccess, ""},
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
