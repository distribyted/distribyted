package iio

import "io"

type Reader interface {
	io.ReaderAt
	io.Closer
	io.Reader
}
