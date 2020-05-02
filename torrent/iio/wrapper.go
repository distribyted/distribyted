package iio

import (
	"io"
	"sync"
)

type readAtWrapper struct {
	mu sync.Mutex

	io.ReadSeeker
	io.ReaderAt
}

func NewReadAtWrapper(r io.ReadSeeker) io.ReaderAt {
	return &readAtWrapper{ReadSeeker: r}
}

func (rw *readAtWrapper) ReadAt(p []byte, off int64) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	_, err := rw.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return rw.Read(p)
}
