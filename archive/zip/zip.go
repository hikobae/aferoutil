package zip

import (
	ziplib "archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

// Unarchive unarchive srcFile to dstDir
func Unarchive(fs afero.Fs, srcFile, dstDir string) (err error) {
	af := afero.Afero{Fs: fs}

	var exists bool
	exists, err = af.DirExists(dstDir)
	if err != nil {
		return
	}
	if exists {
		return fmt.Errorf("Destination directory already exists: %q", dstDir)
	}

	err = af.MkdirAll(dstDir, 0777)
	if err != nil {
		return
	}

	var f afero.File
	f, err = af.Open(srcFile)
	if err != nil {
		return
	}
	defer er(f.Close, &err)

	var fi os.FileInfo
	fi, err = f.Stat()
	if err != nil {
		return
	}

	var r *ziplib.Reader
	r, err = ziplib.NewReader(f, fi.Size())
	if err != nil {
		return
	}

	for _, f := range r.File {
		p := filepath.Join(dstDir, f.Name)

		err = unarchiveFile(af, f, p)
		if err != nil {
			return
		}

		// afero.MemMapFs で io.ReadCloser を Close() すると Chtimes() の結果が消えてしまうため (Bug?),
		// Close() 後に Chtimes() する.
		err = af.Chtimes(p, time.Now(), f.Modified)
		if err != nil {
			return
		}
	}
	return
}

func unarchiveFile(af afero.Afero, f *ziplib.File, p string) (err error) {
	var r io.ReadCloser
	r, err = f.Open()
	if err != nil {
		return
	}
	defer er(r.Close, &err)

	if f.FileInfo().IsDir() {
		af.MkdirAll(p, f.Mode())
	} else {
		var fp afero.File
		fp, err = af.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return
		}
		defer er(fp.Close, &err)

		_, err = io.Copy(fp, r)
		if err != nil {
			return
		}
	}
	return
}

func er(f func() error, oldErr *error) {
	err := f()
	if *oldErr == nil {
		*oldErr = err
	}
}
