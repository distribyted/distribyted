package iio

import (
	"errors"
	"io"
	"io/ioutil"
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

type unbufferedReaderAt struct {
	io.Reader
	N int64
}

func NewUnbufferedReaderAt(r io.Reader) io.ReaderAt {
	return &unbufferedReaderAt{Reader: r}
}

// TODO not working properly
func (u *unbufferedReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	if off < u.N {
		return 0, errors.New("invalid offset")
	}
	diff := off - u.N
	written, err := io.CopyN(ioutil.Discard, u, diff)
	u.N += written
	if err != nil {
		return 0, err
	}

	n, err = u.Read(p)
	u.N += int64(n)
	return
}
