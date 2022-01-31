package vfs

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestSevenZipFS(t *testing.T) {
	require := require.New(t)

	f, err := os.Open("./testdata/sample.7z")
	require.NoError(err)

	fi, err := f.Stat()
	require.NoError(err)

	sevenZip, err := NewSevenZIP(f, fi.Size())
	require.NoError(err)

	err = fstest.TestFS(sevenZip,
		"testdata/quote1.txt",
		"testdata/already-compressed.jpg",
		"testdata/proverbs",
	)

	require.NoError(err)
}
