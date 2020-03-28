package aferoutil

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/spf13/afero"
)

type testcase struct {
	root  string
	dirs  []string
	files []string
}

func TestListFiles(t *testing.T) {
	testcases := []testcase{
		{
			root:  "tmpdir",
			dirs:  []string{"b/", "c/", "c/d/", "p/"},
			files: []string{"a.txt", "z.txt", "b/c.txt", "b/o.ml", "c/y.c"},
		},
		{
			root:  "empty",
			dirs:  []string{},
			files: []string{},
		},
	}

	for _, tc := range testcases {
		testListFilesFile(t, tc)
	}
}

func testListFilesFile(t *testing.T, tc testcase) {
	fs := afero.NewMemMapFs()
	af := afero.Afero{Fs: fs}

	af.MkdirAll(tc.root, 0666)

	for _, file := range tc.files {
		p := filepath.Join(tc.root, file)
		if err := af.WriteFile(p, []byte{}, 0644); err != nil {
			t.Fatalf("Fail to WriteFile(%q): %v", p, err)
		}
	}

	for _, dir := range tc.dirs {
		p := filepath.Join(tc.root, dir)
		if err := af.MkdirAll(p, 0666); err != nil {
			t.Fatalf("Fail to MkdirAll(%q): %v", p, err)
		}
	}

	expected := append(tc.files, tc.dirs...)
	for i, p := range expected {
		expected[i] = filepath.Join(tc.root, p)
	}
	sort.Strings(expected)

	files, err := ListFiles(fs, tc.root)
	if err != nil {
		t.Errorf("expected nil, but %v", err)
		return
	}
	t.Logf("files: %v", files)
	sort.Strings(files)

	if !reflect.DeepEqual(expected, files) {
		t.Errorf("not equals: expected %v, but %v", expected, files)
	}
}
