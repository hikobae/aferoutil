package zip

import (
	ziplib "archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/afero"
)

type archiveTest struct {
	zipPath string
	srcDir  string
	files   []archiveTestFile
}

type archiveTestFile struct {
	name     string
	modified time.Time
	content  []byte
	mode     os.FileMode
}

var archiveTests = []archiveTest{
	{
		zipPath: "a.zip",
		srcDir:  "src",
		files: []archiveTestFile{
			{
				name:     "a.txt",
				modified: time.Date(2020, 3, 30, 8, 14, 37, 318562300, time.FixedZone("JST", 9*60*60)),
				content:  []byte("this is test"),
			},
		},
	},
}

func TestArchive(t *testing.T) {
	fs := afero.NewMemMapFs()
	for _, tc := range archiveTests {
		testArchiveFile(t, fs, tc)
	}
}

func findArchiveTest(name string, files []archiveTestFile) *archiveTestFile {
	for _, f := range files {
		if name == f.name {
			return &f
		}
	}
	return nil
}

func testArchiveFile(t *testing.T, fs afero.Fs, tc archiveTest) {
	af := afero.Afero{Fs: fs}

	// prepare input files
	if err := af.MkdirAll(tc.srcDir, 0666); err != nil {
		t.Fatalf("fail to MkdirAll(%q): %v", tc.srcDir, err)
	}
	defer func() {
		if err := af.RemoveAll(tc.srcDir); err != nil {
			t.Errorf("expected nil, but %v", err)
		}
	}()

	for _, file := range tc.files {
		p := filepath.Join(tc.srcDir, file.name)
		if err := af.WriteFile(p, file.content, file.mode); err != nil {
			t.Fatalf("fail to WriteFile(%q): %v", p, err)
		}
	}

	if err := Archive(fs, tc.srcDir, tc.zipPath); err != nil {
		t.Fatalf("fail to Archive(): %v", err)
	}

	f, err := af.Open(tc.zipPath)
	if err != nil {
		t.Fatalf("fail to Open(%q): %v", tc.zipPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("fail to Close(%q): %v", tc.zipPath, err)
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("fail to Stat(%q): %v", tc.zipPath, err)
	}

	r, err := ziplib.NewReader(f, fi.Size())
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range r.File {
		testArchivedFile(t, fs, tc, file)
	}
}

func testArchivedFile(t *testing.T, fs afero.Fs, tc archiveTest, file *ziplib.File) {
	f, err := file.Open()
	if err != nil {
		t.Errorf("expected nil, but %v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			t.Errorf("expected nil, but %v", err)
		}
	}()

	t.Log(file.Name)
	if file.FileInfo().IsDir() {
		return
	}

	found := findArchiveTest(file.Name, tc.files)
	if found == nil {
		t.Errorf("not found %q in %v", file.Name, tc.files)
		return
	}

	buf := make([]byte, file.UncompressedSize)
	_, err = io.ReadFull(f, buf)
	if err != nil {
		t.Errorf("expected nil, but %v", err)
		return
	}

	if bytes.Compare(found.content, buf) != 0 {
		t.Errorf("expected %v, but %v", found.content, buf)
		return
	}
}
