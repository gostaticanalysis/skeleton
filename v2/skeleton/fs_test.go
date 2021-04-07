package skeleton_test

import "testing"

func TestCreateDir(t *testing.T) {
	t.Parallel()

	FS := func(files ...string) []string {return files}
	cases := map[string]struct {
		files S
	}{
		"empty": {FS()},
		"single": {FS("a")},
		"multi": {FS("a", "b")},
		"subdir": {FS("a", "b", "c/c")},
		"subdirs": {FS("a", "b", "c/c", "d/d")},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
