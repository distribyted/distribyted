package fuse

import (
	"errors"
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
	fh *fileHandler
}

func NewFS(fss []fs.Filesystem) fuse.FileSystemInterface {
	return &FS{
		fh: &fileHandler{fss: fss},
	}
}

func (fs *FS) Open(path string, flags int) (errc int, fh uint64) {
	fh, err := fs.fh.OpenHolder(path)
	if err == os.ErrNotExist {
		logrus.WithField("path", path).Warn("file does not exists")
		return -fuse.ENOENT, fhNone

	}
	if err != nil {
		logrus.WithError(err).WithField("path", path).Error("error opening file")
		return -fuse.EIO, fhNone
	}

	return 0, fh
}

func (fs *FS) Opendir(path string) (errc int, fh uint64) {
	return fs.Open(path, 0)
}

func (cfs *FS) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	if path == "/" {
		stat.Mode = fuse.S_IFDIR | 0555
		return 0
	}

	file, err := cfs.fh.GetFile(path, fh)
	if err == os.ErrNotExist {
		logrus.WithField("path", path).Warn("file does not exists")
		return -fuse.ENOENT

	}
	if err != nil {
		logrus.WithError(err).WithField("path", path).Error("error getting holder when reading file attributes")
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
	file, err := fs.fh.GetFile(path, fh)
	if err == os.ErrNotExist {
		logrus.WithField("path", path).Error("file not found on READ operation")
		return -fuse.ENOENT

	}
	if err != nil {
		logrus.WithError(err).WithField("path", path).Error("error getting holder reading data from file")
		return -fuse.EIO
	}

	end := int(math.Min(float64(len(dest)), float64(int64(file.Size())-off)))
	if end < 0 {
		end = 0
	}

	buf := dest[:end]

	n, err := file.ReadAt(buf, off)
	if err != nil && err != io.EOF {
		logrus.WithError(err).WithField("path", path).Error("error reading data")
		return -fuse.EIO
	}

	dest = buf[:n]

	return n
}

func (fs *FS) Release(path string, fh uint64) (errc int) {
	if err := fs.fh.Remove(fh); err != nil {
		logrus.WithError(err).WithField("path", path).Error("error getting holder when releasing file")
		return -fuse.EIO
	}

	return 0
}

func (fs *FS) Releasedir(path string, fh uint64) int {
	return fs.Release(path, fh)
}

func (fs *FS) Readdir(path string,
	fill func(name string, stat *fuse.Stat_t, ofst int64) bool,
	ofst int64,
	fh uint64) (errc int) {
	fill(".", nil, 0)
	fill("..", nil, 0)

	//TODO improve this function to make use of fh index if possible
	paths, err := fs.fh.ListDir(path)
	if err != nil {
		logrus.WithField("path", path).Error("error reading directory")
		return -fuse.EIO
	}

	for _, p := range paths {
		if !fill(p, nil, 0) {
			logrus.WithField("path", p).Error("error adding directory")
			break
		}
	}

	return 0
}

const fhNone = ^uint64(0)

var ErrHolderEmpty = errors.New("file holder is empty")
var ErrBadHolderIndex = errors.New("holder index too big")

type fileHandler struct {
	mu     sync.Mutex
	opened []fs.File
	fss    []fs.Filesystem
}

func (fh *fileHandler) GetFile(path string, fhi uint64) (fs.File, error) {
	fh.mu.Lock()
	defer fh.mu.Unlock()

	if fhi == fhNone {
		return fh.lookupFile(path)
	}

	return fh.get(fhi)
}

func (fh *fileHandler) ListDir(path string) ([]string, error) {
	fh.mu.Lock()
	defer fh.mu.Unlock()

	var out []string
	for _, ifs := range fh.fss {
		files, err := ifs.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for p := range files {
			out = append(out, p)
		}
	}

	return out, nil
}

func (fh *fileHandler) OpenHolder(path string) (uint64, error) {
	fh.mu.Lock()
	defer fh.mu.Unlock()

	file, err := fh.lookupFile(path)
	if err != nil {
		return fhNone, err
	}

	for i, old := range fh.opened {
		if old == nil {
			fh.opened[i] = file
			return uint64(i), nil
		}
	}
	fh.opened = append(fh.opened, file)

	return uint64(len(fh.opened) - 1), nil
}

func (fh *fileHandler) get(fhi uint64) (fs.File, error) {
	if int(fhi) >= len(fh.opened) {
		return nil, ErrBadHolderIndex
	}
	// TODO check opened slice to avoid panics
	h := fh.opened[int(fhi)]
	if h == nil {
		return nil, ErrHolderEmpty
	}

	return h, nil
}

func (fh *fileHandler) Remove(fhi uint64) error {
	fh.mu.Lock()
	defer fh.mu.Unlock()

	if fhi == fhNone {
		return nil
	}

	// TODO check opened slice to avoid panics
	f := fh.opened[int(fhi)]
	if f == nil {
		return ErrHolderEmpty
	}

	if err := f.Close(); err != nil {
		return err
	}

	fh.opened[int(fhi)] = nil

	return nil
}

func (fh *fileHandler) lookupFile(path string) (fs.File, error) {
	for _, f := range fh.fss {
		file, err := f.Open(path)
		if err == os.ErrNotExist {
			continue
		}
		if err != nil {
			return nil, err
		}

		if file != nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}
