package main

import (
	"os"

	"github.com/hikobae/aferoutil/archive/zip"
	"github.com/spf13/afero"
)

func main() {
	fs := afero.NewOsFs()
	if err := zip.Archive(fs, os.Args[1], os.Args[2]); err != nil {
		panic(err)
	}
}
