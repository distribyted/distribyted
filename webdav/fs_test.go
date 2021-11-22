package webdav

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/distribyted/distribyted/fs"
	"github.com/stretchr/testify/require"
)

func TestWebDAVFilesystem(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	mfs := fs.NewMemory()
	mf := fs.NewMemoryFile([]byte("test file content."))
	err := mfs.Storage.Add(mf, "/folder/file.txt")
	require.NoError(err)

	wfs := newFS(mfs)

	dir, err := wfs.OpenFile(context.Background(), "/", 0, 0)
	require.NoError(err)

	fi, err := dir.Readdir(0)
	require.NoError(err)
	require.Len(fi, 1)
	require.Equal("folder", fi[0].Name())

	file, err := wfs.OpenFile(context.Background(), "/folder/file.txt", 0, 0)
	require.NoError(err)
	_, err = file.Readdir(0)
	require.ErrorIs(err, os.ErrInvalid)

	n, err := file.Seek(5, io.SeekStart)
	require.NoError(err)
	require.Equal(int64(5), n)

	br := make([]byte, 4)
	nn, err := file.Read(br)
	require.NoError(err)
	require.Equal(4, nn)
	require.Equal([]byte("file"), br)

	n, err = file.Seek(0, io.SeekStart)
	require.NoError(err)
	require.Equal(int64(0), n)

	nn, err = file.Read(br)
	require.NoError(err)
	require.Equal(4, nn)
	require.Equal([]byte("test"), br)
}
