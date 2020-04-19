package zip

import (
	"archive/zip"
	ziplib "archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

// Archive archive files in srcDir to dstFile
//
// dstDir の構成が以下のとき, dstFile (zip) に以下のように圧縮する.
//
// dstDir/
// |- x.txt
// :
//
// dstFile (zip)
// |- x.txt
// :
func Archive(fs afero.Fs, srcDir, dstFile string) (err error) {
	af := afero.Afero{Fs: fs}

	var exists bool
	exists, err = af.Exists(dstFile)
	if err != nil {
		return
	}
	if exists {
		return fmt.Errorf("Destination file already exists: %q", dstFile)
	}

	var f afero.File
	f, err = af.Create(dstFile)
	if err != nil {
		return
	}
	defer er(f.Close, &err)

	w := zip.NewWriter(f)
	defer er(w.Close, &err)

	af.Walk(srcDir, func(path string, info os.FileInfo, incomingErr error) (err error) {
		if incomingErr != nil {
			err = incomingErr
			return
		}

		var header *ziplib.FileHeader
		header, err = zip.FileInfoHeader(info)
		if err != nil {
			return
		}

		var rel string
		rel, err = filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." || rel == ".." {
			return
		}

		header.Name = filepath.ToSlash(rel)
		header.Modified = info.ModTime().In(time.Local)

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		var h io.Writer
		h, err = w.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		var f afero.File
		f, err = af.Open(path)
		if err != nil {
			return
		}
		defer er(f.Close, &err)

		_, err = io.Copy(h, f)
		return
	})
	return
}
