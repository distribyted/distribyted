package vfs

import (
	"errors"
	"io"
	"io/fs"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/bodgit/sevenzip"
)

var SevenZipFactory = func(f fs.File) (fs.FS, error) {
	ra, ok := f.(io.ReaderAt)
	if !ok {
		return nil, errors.New("ReadAt is needed")
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return NewSevenZIP(ra, fi.Size())
}

var _ fs.FS = &SevenZIP{}

type SevenZIP struct {
	r *sevenzip.Reader

	onceFileList sync.Once
	fil          *FileInfoList
}

func NewSevenZIP(r io.ReaderAt, size int64) (*SevenZIP, error) {
	zr, err := sevenzip.NewReader(r, size)
	return &SevenZIP{r: zr,
		fil: NewFileInfoList(),
	}, err
}

func (vfs *SevenZIP) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	vfs.load()

	// special case for root
	if name == "." {
		return newDir(newDirFileInfo(name), fileInfosToDirEntries(vfs.fil.Lookup(name))), nil
	}

	for _, f := range vfs.r.File {
		fi := newSevenZIPFileInfo(f.FileHeader)

		if fi.Path() == name && !fi.IsDir() {
			fr, err := f.Open()
			if err != nil {
				return nil, err
			}

			return newFile(fi, &fixZeroBytesReader{fr}), nil
		}

		if fi.Path() == name && fi.IsDir() {
			return newDir(fi, fileInfosToDirEntries(vfs.fil.Lookup(name))), nil
		}
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

func (vfs *SevenZIP) load() error {
	vfs.onceFileList.Do(
		func() {
			for _, f := range vfs.r.File {
				vfs.fil.Add(newSevenZIPFileInfo(f.FileHeader))
			}

			vfs.fil.Sort()
		})

	return nil
}

var _ FileInfoPath = &sevenZIPFileInfo{}

func newSevenZIPFileInfo(fh sevenzip.FileHeader) *sevenZIPFileInfo {
	return &sevenZIPFileInfo{
		fh:   fh,
		name: strings.TrimSuffix(fh.Name, "/"),
	}
}

type sevenZIPFileInfo struct {
	fh   sevenzip.FileHeader
	name string
}

func (fi *sevenZIPFileInfo) Path() string {
	return fi.name
}

func (fi *sevenZIPFileInfo) Name() string {
	base := path.Base(fi.name)
	return base
}
func (fi *sevenZIPFileInfo) Size() int64 {
	return int64(fi.fh.UncompressedSize)
}
func (fi *sevenZIPFileInfo) Mode() fs.FileMode {
	return fi.fh.Mode()
}
func (fi *sevenZIPFileInfo) ModTime() time.Time {
	return fi.fh.Modified.UTC()
}
func (fi *sevenZIPFileInfo) IsDir() bool {
	return fi.Mode().IsDir()
}
func (fi *sevenZIPFileInfo) Sys() interface{} {
	return nil
}

// used to fix an io.EOF error when the byte slice is nil
type fixZeroBytesReader struct {
	io.ReadCloser
}

func (r *fixZeroBytesReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)

	if len(p) == 0 && err == io.EOF {
		err = nil
	}

	return
}
