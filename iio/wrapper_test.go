package iio_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/distribyted/distribyted/fs"
	"github.com/distribyted/distribyted/iio"
	"github.com/stretchr/testify/require"
)

var testData []byte = []byte("Hello World")

func TestReadAtWrapper(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	br := bytes.NewReader(testData)

	r := iio.NewReadAtWrapper(br)
	defer r.Close()

	toRead := make([]byte, 5)
	n, err := r.ReadAt(toRead, 6)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal("World", string(toRead))

	n, err = r.ReadAt(toRead, 0)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal("Hello", string(toRead))
}

func TestSeekerWrapper(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	mf := fs.NewMemoryFile(testData)

	r := iio.NewSeekerWrapper(mf, mf.Size())
	defer r.Close()

	n, err := r.Seek(6, io.SeekStart)
	require.NoError(err)
	require.Equal(int64(6), n)

	toRead := make([]byte, 5)
	nn, err := r.Read(toRead)
	require.NoError(err)
	require.Equal(5, nn)
	require.Equal("World", string(toRead))
}
