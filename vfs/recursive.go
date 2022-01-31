package vfs

import (
	"io/fs"
	"path"
	"strings"
)

type FSFactory func(f fs.File) (fs.FS, error)

type Recursive struct {
	root      fs.FS
	factories map[string]FSFactory
}

func NewRecursive(root fs.FS, fsFactories map[string]FSFactory) *Recursive {
	return &Recursive{
		root:      root,
		factories: fsFactories,
	}
}

func (vfs *Recursive) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	elms, err := vfs.parseElements(name)
	if err != nil {
		return nil, err
	}

	cfs := vfs.root
	key := ""
	var file fs.File
	for _, elm := range elms {
		f, err := elm.fs.Open(elm.key)
		if err != nil {
			return nil, err
		}

		file = f
		key = elm.key
		cfs = elm.fs
	}

	if file == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	fs, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if fs.IsDir() {
		return file, nil
	}

	return cfs.Open(key)
}

func (vfs *Recursive) parseElements(name string) ([]element, error) {
	elementKeys := strings.Split(name, "/")

	if len(elementKeys) == 0 {
		return []element{{vfs.root, "."}}, nil
	}

	key := "."
	actualFS := vfs.root
	var elements []element
	for _, ekey := range elementKeys {
		key = path.Join(key, ekey)
		info, err := fs.Stat(actualFS, key)
		if err != nil {
			return nil, err
		}

		elements = append(elements, element{actualFS, key})

		if info.IsDir() {
			continue
		}

		ext := path.Ext(key)
		factory, ok := vfs.factories[ext]
		if !ok {
			continue
		}

		f, err := actualFS.Open(key)
		if err != nil {
			return nil, err
		}

		childFS, err := factory(f)
		if err != nil {
			return nil, err
		}

		key = "."
		actualFS = childFS

	}

	return elements, nil
}

type element struct {
	fs  fs.FS
	key string
}
