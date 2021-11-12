package skeleton_test

import (
	"bytes"
	"flag"
	"io"
	"io/fs"
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
	os.Setenv("SKELETON_PREFIX", "")
}

func TestSkeletonRun(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		dirinit string
		args    string
		input   string

		wantExitCode int
		wantOutput   string
	}{
		"nooption":              {"", "example.com/example", "", skeleton.ExitSuccess, ""},
		"overwrite-cancel":      {"-- example/go.mod --\n// empty", "example.com/example", "1\n", skeleton.ExitSuccess, ""},
		"overwrite-force":       {"-- example/go.mod --\n// empty", "example.com/example", "2\n", skeleton.ExitSuccess, ""},
		"overwrite-confirm-yes": {"-- example/go.mod --\n// empty", "example.com/example", "3\ny\n", skeleton.ExitSuccess, ""},
		"overwrite-confirm-no":  {"-- example/go.mod --\n// empty", "example.com/example", "3\nn\n", skeleton.ExitSuccess, ""},
		"overwrite-newonly":     {"-- example/go.mod --\n// empty", "example.com/example", "4\n", skeleton.ExitSuccess, ""},
		"plugin":                {"", "-plugin example.com/example", "", skeleton.ExitSuccess, ""},
		"nocmd":                 {"", "-cmd=false example.com/example", "", skeleton.ExitSuccess, ""},
		"onlypkgname":           {"", "example", "", skeleton.ExitSuccess, ""},
		"version":               {"", "-v", "", skeleton.ExitSuccess, "skeleton version\n"},
	}

	if flagUpdate {
		removeGolden(t)
	}

	for name, tt := range cases {
		name, tt := name, tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			dir := t.TempDir()
			if tt.dirinit != "" {
				dirInit(t, dir, tt.dirinit)
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

			got := txtarDir(t, s.Dir)

			if flagUpdate {
				writeTestData(t, name, got)
				return
			}

			want := readTestData(t, name)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func dirInit(t *testing.T, root, s string) {
	t.Helper()
	fsys := txtarfs.As(txtar.Parse([]byte(s)))
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) (rerr error) {
		if err != nil {
			return err
		}

		// directory would create with a file
		if d.IsDir() {
			return nil
		}

		dstPath := filepath.Join(root, filepath.FromSlash(path))

		src, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		fi, err := src.Stat()
		if err != nil {
			return err
		}

		if fi.Size() == 0 {
			return nil
		}

		err = os.MkdirAll(filepath.Dir(dstPath), 0700)
		if err != nil {
			return err
		}

		dst, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer func() {
			if err := dst.Close(); err != nil && rerr == nil {
				rerr = err
			}
		}()

		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal("unexpected error", err)
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
	f, err := os.Create(filepath.Join("testdata", name+".golden"))
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
	got, err := os.ReadFile(filepath.Join("testdata", name+".golden"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	return string(got)
}

func removeGolden(t *testing.T) {
	t.Helper()
	err := filepath.Walk("testdata", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(path) != ".golden" {
			return nil
		}

		if err := os.Remove(path); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal("unexpected error", err)
	}
}
