package iio

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

var testData []byte = []byte("Hello World")

func TestReadData(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	br := bytes.NewReader(testData)
	r, err := NewDiskTeeReader(br)
	require.NoError(err)

	toRead := make([]byte, 5)

	n, err := r.ReadAt(toRead, 6)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal("World", string(toRead))

	r.ReadAt(toRead, 0)
	require.NoError(err)
	require.Equal(5, n)
	require.Equal("Hello", string(toRead))
}

func TestReadDataEOF(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	br := bytes.NewReader(testData)
	r, err := NewDiskTeeReader(br)
	require.NoError(err)

	toRead := make([]byte, 6)

	n, err := r.ReadAt(toRead, 6)
	require.Equal(io.EOF, err)
	require.Equal(5, n)
	require.Equal("World\x00", string(toRead))
}
