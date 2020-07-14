package distribyted

import (
	"net/http"
	"path"
	"strings"
)

type binaryFileSystem struct {
	fs   http.FileSystem
	base string
}

func NewBinaryFileSystem(fs http.FileSystem, base string) *binaryFileSystem {
	return &binaryFileSystem{
		fs:   fs,
		base: base,
	}
}

func (fs *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func (fs *binaryFileSystem) Open(name string) (http.File, error) {
	return fs.fs.Open(path.Join(fs.base, name))
}
