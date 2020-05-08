package distribyted

import (
	"net/http"
	"strings"
)

type binaryFileSystem struct {
	http.FileSystem
}

func NewBinaryFileSystem() *binaryFileSystem {
	return &binaryFileSystem{Assets}
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
