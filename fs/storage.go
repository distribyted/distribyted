package fs

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

const separator = "/"

type FsFactory func(f File) (Filesystem, error)

var SupportedFactories = map[string]FsFactory{
	".zip": func(f File) (Filesystem, error) {
		return NewZip(f, f.Size()), nil
	},
}

type storage struct {
	factories map[string]FsFactory

	files       map[string]File
	filesystems map[string]Filesystem
	children    map[string]map[string]File
}

func newStorage(factories map[string]FsFactory) *storage {
	return &storage{
		files:       make(map[string]File, 0),
		children:    make(map[string]map[string]File, 0),
		filesystems: make(map[string]Filesystem, 0),
		factories:   factories,
	}
}

func (s *storage) Has(path string) bool {
	path = clean(path)

	f := s.files[path]
	if f != nil {
		return true
	}

	if f, _ := s.getFileFromFs(path); f != nil {
		return true
	}

	return false
}

func (s *storage) Add(f File, p string) error {
	p = clean(p)
	if s.Has(p) {
		if dir, err := s.Get(p); err == nil {
			if !dir.IsDir() {
				return os.ErrExist
			}
		}

		return nil
	}

	ext := path.Ext(p)
	if ffs := s.factories[ext]; ffs != nil {
		fs, err := ffs(f)
		if err != nil {
			return err
		}

		s.filesystems[p] = fs
	} else {
		s.files[p] = f
	}

	s.createParent(p, f)

	return nil
}

func (s *storage) createParent(path string, f File) error {
	base, filename := filepath.Split(path)
	base = clean(base)

	if err := s.Add(&Dir{}, base); err != nil {
		return err
	}

	if _, ok := s.children[base]; !ok {
		s.children[base] = make(map[string]File, 0)
	}

	if filename != "" {
		s.children[base][filename] = f
	}

	return nil
}

func (s *storage) Children(path string) map[string]File {
	path = clean(path)

	out, err := s.getDirFromFs(path)
	if err == nil {
		return out
	}

	l := make(map[string]File, 0)
	for n, f := range s.children[path] {
		l[n] = f
	}

	return l
}

func (s *storage) Get(path string) (File, error) {
	path = clean(path)
	if !s.Has(path) {
		return nil, os.ErrNotExist
	}

	file, ok := s.files[path]
	if ok {
		return file, nil
	}

	return s.getFileFromFs(path)
}

func (s *storage) getFileFromFs(p string) (File, error) {
	for fsp, fs := range s.filesystems {
		if strings.HasPrefix(p, fsp) {
			return fs.Open(separator + strings.TrimPrefix(p, fsp))
		}
	}

	return nil, os.ErrNotExist
}

func (s *storage) getDirFromFs(p string) (map[string]File, error) {
	for fsp, fs := range s.filesystems {
		if strings.HasPrefix(p, fsp) {
			path := strings.TrimPrefix(p, fsp)
			return fs.ReadDir(path)
		}
	}

	return nil, os.ErrNotExist
}

func clean(path string) string {
	return filepath.Clean(separator + strings.ReplaceAll(path, "\\", "/"))
}
