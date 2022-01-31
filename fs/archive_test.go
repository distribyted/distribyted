package fs

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/distribyted/distribyted/iio"
	"github.com/stretchr/testify/require"
)

var fileContent []byte = []byte("Hello World")

func TestRarFilesystem(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	f, err := os.Open("./testdata/sample.rar")
	require.NoError(err)

	fs, err := f.Stat()
	require.NoError(err)

	rfs := NewArchive(f, fs.Size(), &Rar{})

	f1, err := rfs.Open("/testdata/quote1.txt")
	require.NoError(err)
	require.Equal(int64(58), f1.Size())

	f2, err := rfs.Open("/testdata/proverbs/proverb1.txt")
	require.NoError(err)
	require.Equal(int64(54), f2.Size())

	b := make([]byte, 5)

	n, err := f1.Read(b)
	require.NoError(err)
	require.Equal(int64(3), n)
	require.Equal("tes", string(b))

	n, err = f2.Read(b)
	require.NoError(err)
	require.Equal(int64(3), n)
	require.Equal("tes", string(b))
}

func TestZipFilesystem(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	zReader, len := createTestZip(require)

	zfs := NewArchive(zReader, len, &Zip{})

	files, err := zfs.ReadDir("/path/to/test/file")
	require.NoError(err)

	require.Len(files, 1)
	f := files["1.txt"]
	require.NotNil(f)

	out := make([]byte, 11)
	n, err := f.Read(out)
	require.Equal(io.EOF, err)
	require.Equal(11, n)
	require.Equal(fileContent, out)

}

func createTestZip(require *require.Assertions) (iio.Reader, int64) {
	buf := bytes.NewBuffer([]byte{})

	zWriter := zip.NewWriter(buf)

	f1, err := zWriter.Create("path/to/test/file/1.txt")
	require.NoError(err)
	_, err = f1.Write(fileContent)
	require.NoError(err)

	err = zWriter.Close()
	require.NoError(err)

	return newCBR(buf.Bytes()), int64(buf.Len())
}

type closeableByteReader struct {
	*bytes.Reader
}

func newCBR(b []byte) *closeableByteReader {
	return &closeableByteReader{
		Reader: bytes.NewReader(b),
	}
}

func (*closeableByteReader) Close() error {
	return nil
}
