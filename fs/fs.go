package fs

import (
	"github.com/ajnavarro/distribyted/iio"
)

type File interface {
	IsDir() bool
	Size() int64

	iio.Reader
}

type Filesystem interface {
	// Open opens the named file for reading. If successful, methods on the
	// returned file can be used for reading; the associated file descriptor has
	// mode O_RDONLY.
	Open(filename string) (File, error)

	// ReadDir reads the directory named by dirname and returns a list of
	// directory entries.
	ReadDir(path string) (map[string]File, error)
}
