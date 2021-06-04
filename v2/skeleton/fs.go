package skeleton

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type overwritePolicy int

const (
	cancel overwritePolicy = iota
	forceOverwrite
	confirm
	newOnly
)

// CreateDir creates files and directries which structure same with the given file system.
// The path of created root directory become the parameter root.
func CreateDir(prompt *Prompt, root string, fsys fs.FS) error {
	policy, err := choosePolicy(prompt, root)
	if err != nil {
		return err
	}

	if policy == forceOverwrite {
		if err := os.RemoveAll(root); err != nil {
			return err
		}
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) (rerr error) {
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

		dst, err := create(prompt, dstPath, policy)
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
		return fmt.Errorf("CreateDir: %w", err)
	}
	return nil
}

func choosePolicy(prompt *Prompt, dir string) (overwritePolicy, error) {
	exist, err := isExist(dir)
	if err != nil {
		return 0, err
	}

	if !exist {
		return cancel, nil
	}

	desc := fmt.Sprintf("%s is already exist, overwrite?", dir)
	opts := []string{
		cancel:         "No (Exit)",
		forceOverwrite: "Remove and create new directory",
		confirm:        "Overwrite existing files with confirmation",
		newOnly:        "Create new files only",
	}
	n, err := prompt.Choose(desc, opts, ">")
	if err != nil {
		return 0, err
	}

	return overwritePolicy(n), nil
}

func isExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func create(prompt *Prompt, path string, policy overwritePolicy) (io.WriteCloser, error) {
	var nopWriter = struct {
		io.Writer
		io.Closer
	}{io.Discard, io.NopCloser(nil)}

	exist, err := isExist(path)
	if err != nil {
		return nil, err
	}	

	if !exist {
		return os.Create(path)
	}

	if policy != confirm {
		return nopWriter, nil
	}

	desc := fmt.Sprintf("%s is already exist, overwrite?", path)
	yesno, err := prompt.YesNo(desc, false, '>')
	if err != nil {
		return nil, err
	}

	if !yesno {
		return nopWriter, nil
	}

	return os.Create(path)
}
