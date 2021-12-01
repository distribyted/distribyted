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

type seekerWrapper struct {
	mu   sync.Mutex
	pos  int64
	size int64

	io.Seeker
	Reader
}

func NewSeekerWrapper(r Reader, size int64) *seekerWrapper {
	return &seekerWrapper{Reader: r}
}

func (r *seekerWrapper) Seek(offset int64, whence int) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch whence {
	case io.SeekStart:
		r.pos = offset
	case io.SeekCurrent:
		r.pos = r.pos + offset
	case io.SeekEnd:
		r.pos = r.size + offset
	}

	return r.pos, nil
}

func (r *seekerWrapper) Read(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n, err := r.ReadAt(p, r.pos)
	r.pos += int64(n)

	return n, err
}
