package aferoutil

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// ListFiles list files and directories recursively.
func ListFiles(fs afero.Fs, root string) ([]string, error) {
	root = filepath.Clean(root)

	files := []string{}
	err := afero.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == root {
			return nil
		}

		files = append(files, path)
		return nil
	})
	return files, err
}
