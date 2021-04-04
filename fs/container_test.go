package fs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainer(t *testing.T) {
	require := require.New(t)

	fss := map[string]Filesystem{
		"/test": &DummyFs{},
	}

	c, err := NewContainerFs(fss)
	require.NoError(err)

	f, err := c.Open("/test/dir/here")
	require.NoError(err)
	require.NotNil(f)

	files, err := c.ReadDir("/")
	require.NoError(err)
	require.Len(files, 1)
}
