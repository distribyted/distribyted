package vfs

import (
	"io/fs"
	"os"
	"path"
)

var _ fs.FS = &Dir{}

type Dir struct {
	path string
}

func NewDir(path string) *Dir {
	return &Dir{
		path: path,
	}
}

func (vfs *Dir) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	return os.Open(path.Join(vfs.path, name))
}
