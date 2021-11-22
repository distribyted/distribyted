package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemory(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	mem := NewMemory()

	mem.Storage.Add(NewMemoryFile([]byte("Hello")), "/dir/here")

	fss := map[string]Filesystem{
		"/test": mem,
	}

	c, err := NewContainerFs(fss)
	require.NoError(err)

	f, err := c.Open("/test/dir/here")
	require.NoError(err)
	require.NotNil(f)
	require.Equal(int64(5), f.Size())
	require.NoError(f.Close())

	files, err := c.ReadDir("/")
	require.NoError(err)
	require.Len(files, 1)

	files, err = c.ReadDir("/test")
	require.NoError(err)
	require.Len(files, 1)

}
