package vfs

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestRecursiveFS(t *testing.T) {
	require := require.New(t)

	dirFS := NewDir("./testdata")

	var factories = map[string]FSFactory{
		".zip": ZipFactory,
		".rar": RarFactory,
		".7z":  SevenZipFactory,
	}

	recursive := NewRecursive(dirFS, factories)

	err := fstest.TestFS(recursive,
		"sample.7z/testdata/quote1.txt",
		"sample.7z/testdata/already-compressed.jpg",
		"sample.7z/testdata/proverbs",
		"sample.rar/testdata/quote1.txt",
		"sample.rar/testdata/already-compressed.jpg",
		"sample.rar/testdata/proverbs",

		"other.zip/sample.7z/testdata/quote1.txt",
		"other.zip/sample.7z/testdata/already-compressed.jpg",
		"other.zip/sample.7z/testdata/proverbs",
		"other.zip/sample.rar/testdata/quote1.txt",
		"other.zip/sample.rar/testdata/already-compressed.jpg",
		"other.zip/sample.rar/testdata/proverbs",
	)

	require.NoError(err)
}
