package iio

import (
	"io"
	"log"
)

func CloseIfCloseable(r interface{}) error {
	log.Println("closing file...")
	if r == nil {
		return nil
	}

	closer, ok := r.(io.Closer)
	if !ok {
		log.Println("file is not implementing close method")
		return nil
	}

	return closer.Close()
}
