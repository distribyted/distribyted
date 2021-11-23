package fs

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileinfo(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	fi := NewFileInfo("name", 42, false)

	require.Equal(fi.IsDir(), false)
	require.Equal(fi.Name(), "name")
	require.Equal(fi.Size(), int64(42))
	require.NotNil(fi.ModTime())
	require.Equal(fi.Mode(), fs.FileMode(0555))
	require.Equal(fi.Sys(), nil)

}
