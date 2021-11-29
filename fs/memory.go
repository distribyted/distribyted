package fs

import (
	"bytes"
)

var _ Filesystem = &Memory{}

type Memory struct {
	Storage *storage
}

func NewMemory() *Memory {
	return &Memory{
		Storage: newStorage(nil),
	}
}

func (fs *Memory) Open(filename string) (File, error) {
	return fs.Storage.Get(filename)
}

func (fs *Memory) ReadDir(path string) (map[string]File, error) {
	return fs.Storage.Children(path)
}

var _ File = &MemoryFile{}

type MemoryFile struct {
	*bytes.Reader
}

func NewMemoryFile(data []byte) *MemoryFile {
	return &MemoryFile{
		Reader: bytes.NewReader(data),
	}
}

func (d *MemoryFile) Size() int64 {
	return int64(d.Reader.Len())
}

func (d *MemoryFile) IsDir() bool {
	return false
}

func (d *MemoryFile) Close() (err error) {
	return
}
