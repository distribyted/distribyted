package vfs

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"path"
	"time"

	"github.com/nwaples/rardecode/v2"
)

var RarFactory = func(f fs.File) (fs.FS, error) {
	ra, ok := f.(io.ReaderAt)
	if !ok {
		return nil, errors.New("ReadAt is needed")
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return NewRar(ra, fi.Size()), nil
}

var _ fs.FS = &Rar{}

type Rar struct {
	r    io.ReaderAt
	size int64
}

func NewRar(r io.ReaderAt, size int64) *Rar {
	// TODO parametrizable buffer size
	// TODO parametrizable password

	return &Rar{r: r, size: size}
}

func (vfs *Rar) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	sr := io.NewSectionReader(vfs.r, 0, vfs.size)
	rr, err := rardecode.NewReader(sr)
	if err != nil {
		return nil, err
	}

	var dirFileInfo fs.FileInfo
	fis := NewFileInfoList()
	for {
		fh, err := rr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		fi := newRarFileInfo(fh)

		if isSymlink(fi) {
			//skip links for now
			continue
		}

		fis.Add(fi)

		if fh.Name == name && !fh.IsDir {
			return newFile(fi, ioutil.NopCloser(rr)), nil
		}

		if fh.Name == name && fh.IsDir {
			dirFileInfo = fi
		}
	}

	// special case for root
	if name == "." {
		dirFileInfo = newDirFileInfo(name)
	}

	// we want to open a folder, so we need to generate it
	if dirFileInfo != nil {
		return newDir(dirFileInfo, fileInfosToDirEntries(fis.Lookup(name))), nil
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

var _ FileInfoPath = &rarFileInfo{}

func newRarFileInfo(fh *rardecode.FileHeader) *rarFileInfo {
	return &rarFileInfo{fh: fh}
}

type rarFileInfo struct {
	fh *rardecode.FileHeader
}

func (fi *rarFileInfo) Path() string {
	return fi.fh.Name
}

func (fi *rarFileInfo) Name() string {
	base := path.Base(fi.fh.Name)
	return base
}
func (fi *rarFileInfo) Size() int64 {
	return fi.fh.UnPackedSize
}
func (fi *rarFileInfo) Mode() fs.FileMode {
	return fi.fh.Mode()
}
func (fi *rarFileInfo) ModTime() time.Time {
	return fi.fh.ModificationTime.UTC()
}
func (fi *rarFileInfo) IsDir() bool {
	return fi.fh.IsDir
}
func (fi *rarFileInfo) Sys() interface{} {
	return nil
}
