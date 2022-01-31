package vfs

import (
	"os"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestRarFS(t *testing.T) {
	require := require.New(t)

	f, err := os.Open("./testdata/sample.rar")
	require.NoError(err)

	fi, err := f.Stat()
	require.NoError(err)

	rar := NewRar(f, fi.Size())

	err = fstest.TestFS(rar,
		"testdata/quote1.txt",
		"testdata/already-compressed.jpg",
		"testdata/proverbs",
	)

	require.NoError(err)
}
