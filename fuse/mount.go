package fuse

import (
	"io"
	"math"
	"os"
	"sync"

	"github.com/ajnavarro/distribyted/fs"
	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/sirupsen/logrus"
)

type FS struct {
	fuse.FileSystemBase

	lock sync.Mutex
	FSS  []fs.Filesystem
}

func (fs *FS) Open(path string, flags int) (errc int, fh uint64) {
	return 0, 0
}

func (fs *FS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	defer fs.synchronize()()

	if path == "/" {
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	}

	file, err := fs.findFile(path)
	if err != nil {
		logrus.WithField("path", path).WithError(err).Warn("error finding file")
		return -fuse.EIO
	}

	if err == os.ErrNotExist {
		logrus.WithField("path", path).Warn("file does not exists")
		return -fuse.ENOENT

	}

	if err != nil {
		logrus.WithField("path", path).WithError(err).Error("error reading file")
		return -fuse.EIO
	}

	if file.IsDir() {
		stat.Mode = fuse.S_IFDIR | 0555
	} else {
		stat.Mode = fuse.S_IFREG | 0444
		stat.Size = file.Size()
	}

	return 0
}

func (fs *FS) Read(path string, dest []byte, off int64, fh uint64) int {
	defer fs.synchronize()()

	file, err := fs.findFile(path)
	if err == os.ErrNotExist {
		logrus.WithField("path", path).Warn("file does not exists")
		return -fuse.ENOENT

	}
	if err != nil {
		logrus.WithField("path", path).WithError(err).Warn("error finding file")
		return -fuse.EIO
	}

	end := int(math.Min(float64(len(dest)), float64(int64(file.Size())-off)))
	if end < 0 {
		end = 0
	}

	buf := dest[:end]

	n, err := file.ReadAt(buf, off)
	if err != nil && err != io.EOF {
		logrus.WithError(err).Error("error reading data")
		return -fuse.EIO
	}

	dest = buf[:n]

	return n
}

func (fs *FS) Release(path string, fh uint64) (errc int) {
	defer fs.synchronize()()
	return 0
}

func (fs *FS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	defer fs.synchronize()()

	fill(".", nil, 0)
	fill("..", nil, 0)

	for _, ifs := range fs.FSS {
		files, err := ifs.ReadDir(path)
		if err != nil {
			return -fuse.EIO
		}
		for p := range files {
			if !fill(p, nil, 0) {
				logrus.WithField("path", p).Error("error adding directory")
				break
			}
		}
	}

	return 0
}

func (fs *FS) findFile(path string) (fs.File, error) {
	for _, f := range fs.FSS {
		file, err := f.Open(path)
		if err == os.ErrNotExist {
			continue
		}
		if err != nil {
			return nil, err
		}

		// TODO add file to list of opened files to be able to close it when not in use
		if file != nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}

func (fs *FS) synchronize() func() {
	fs.lock.Lock()
	return func() {
		fs.lock.Unlock()
	}
}
