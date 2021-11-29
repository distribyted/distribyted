package iio

import (
	"io"
	"sync"
)

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
