package iio

import (
	"io"

	log "github.com/sirupsen/logrus"
)

func CloseIfCloseable(r interface{}) error {
	log.Debug("closing file...")
	if r == nil {
		return nil
	}

	closer, ok := r.(io.Closer)
	if !ok {
		log.Debug("file is not implementing close method")
		return nil
	}

	return closer.Close()
}
