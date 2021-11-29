package fs

import (
	"os"
	"path"
	"strings"
)

const separator = "/"

type FsFactory func(f File) (Filesystem, error)

var SupportedFactories = map[string]FsFactory{
	".zip": func(f File) (Filesystem, error) {
		return NewArchive(f, f.Size(), &Zip{}), nil
	},
	".rar": func(f File) (Filesystem, error) {
		return NewArchive(f, f.Size(), &Rar{}), nil
	},
	".7z": func(f File) (Filesystem, error) {
		return NewArchive(f, f.Size(), &SevenZip{}), nil
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
		files:       make(map[string]File),
		children:    make(map[string]map[string]File),
		filesystems: make(map[string]Filesystem),
		factories:   factories,
	}
}

func (s *storage) Clear() {
	s.files = make(map[string]File)
	s.children = make(map[string]map[string]File)
	s.filesystems = make(map[string]Filesystem)

	s.Add(&Dir{}, "/")
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

func (s *storage) AddFS(fs Filesystem, p string) error {
	p = clean(p)
	if s.Has(p) {
		if dir, err := s.Get(p); err == nil {
			if !dir.IsDir() {
				return os.ErrExist
			}
		}

		return nil
	}

	s.filesystems[p] = fs
	return s.createParent(p, &Dir{})
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

	return s.createParent(p, f)
}

func (s *storage) createParent(p string, f File) error {
	base, filename := path.Split(p)
	base = clean(base)

	if err := s.Add(&Dir{}, base); err != nil {
		return err
	}

	if _, ok := s.children[base]; !ok {
		s.children[base] = make(map[string]File)
	}

	if filename != "" {
		s.children[base][filename] = f
	}

	return nil
}

func (s *storage) Children(path string) (map[string]File, error) {
	path = clean(path)

	l := make(map[string]File)
	for n, f := range s.children[path] {
		l[n] = f
	}

	if _, ok := s.children[path]; ok {
		return l, nil
	}

	return s.getDirFromFs(path)

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

func clean(p string) string {
	return path.Clean(separator + strings.ReplaceAll(p, "\\", "/"))
}
