package iio

import (
	"io"
	"sync"
)

type readAtWrapper struct {
	mu sync.Mutex

	io.ReadSeeker
	io.ReaderAt
	io.Closer
}

func NewReadAtWrapper(r io.ReadSeeker) Reader {
	return &readAtWrapper{ReadSeeker: r}
}

func (rw *readAtWrapper) ReadAt(p []byte, off int64) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	_, err := rw.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return io.ReadAtLeast(rw, p, len(p))
}

func (rw *readAtWrapper) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	c, ok := rw.ReadSeeker.(io.Closer)
	if !ok {
		return nil
	}

	return c.Close()
}
