package vfs

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
)

var ZipFactory = func(f fs.File) (fs.FS, error) {
	ra, ok := f.(io.ReaderAt)
	if !ok {
		return nil, errors.New("ReadAt is needed")
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return zip.NewReader(ra, fi.Size())
}
