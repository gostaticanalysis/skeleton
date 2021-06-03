package skeleton

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CreateDir creates files and directries which structure same with the given file system.
// The path of created root directory become the parameter root.
func CreateDir(root string, fsys fs.FS) error {
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) (rerr error) {
		if err != nil {
			return err
		}

		dstPath := filepath.Join(root, filepath.FromSlash(path))

		if d.IsDir() {
			return os.MkdirAll(dstPath, 0700)
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

		src, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

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
