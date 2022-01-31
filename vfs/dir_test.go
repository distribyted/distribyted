package vfs

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestDirFS(t *testing.T) {
	require := require.New(t)

	dir := NewDir("./testdata")

	err := fstest.TestFS(dir,
		"sample.rar",
		"sample.7z",
	)

	require.NoError(err)
}
