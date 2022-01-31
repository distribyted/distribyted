package fs

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/bodgit/sevenzip"
	"github.com/distribyted/distribyted/iio"
	"github.com/nwaples/rardecode/v2"
)

var _ loader = &Zip{}

type Zip struct {
}

func (fs *Zip) getFiles(reader iio.Reader, size int64) (map[string]*ArchiveFile, error) {
	zr, err := zip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	out := make(map[string]*ArchiveFile)
	for _, f := range zr.File {
		f := f
		if f.FileInfo().IsDir() {
			continue
		}

		rf := func() (iio.Reader, error) {
			zr, err := f.Open()
			if err != nil {
				return nil, err
			}

			return iio.NewDiskTeeReader(zr)
		}

		n := filepath.Join(string(os.PathSeparator), f.Name)
		af := NewArchiveFile(rf, f.FileInfo().Size())

		out[n] = af
	}

	return out, nil
}

var _ loader = &SevenZip{}

type SevenZip struct {
}

func (fs *SevenZip) getFiles(reader iio.Reader, size int64) (map[string]*ArchiveFile, error) {
	r, err := sevenzip.NewReader(reader, size)
	if err != nil {
		return nil, err
	}

	out := make(map[string]*ArchiveFile)
	for _, f := range r.File {
		f := f
		if f.FileInfo().IsDir() {
			continue
		}

		rf := func() (iio.Reader, error) {
			zr, err := f.Open()
			if err != nil {
				return nil, err
			}

			return iio.NewDiskTeeReader(zr)
		}

		af := NewArchiveFile(rf, f.FileInfo().Size())
		n := filepath.Join(string(os.PathSeparator), f.Name)

		out[n] = af
	}

	return out, nil
}

var _ loader = &Rar{}

type Rar struct {
}

func (fs *Rar) getFiles(reader iio.Reader, size int64) (map[string]*ArchiveFile, error) {
	r, err := rardecode.NewReader(iio.NewSeekerWrapper(reader, size))
	if err != nil {
		return nil, err
	}

	out := make(map[string]*ArchiveFile)
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if header.IsDir {
			continue
		}

		rf := func() (iio.Reader, error) {
			return iio.NewDiskTeeReader(r)
		}

		n := filepath.Join(string(os.PathSeparator), header.Name)

		af := NewArchiveFile(rf, header.UnPackedSize)

		out[n] = af
	}

	return out, nil
}

type loader interface {
	getFiles(r iio.Reader, size int64) (map[string]*ArchiveFile, error)
}

var _ Filesystem = &archive{}

type archive struct {
	r iio.Reader
	s *storage

	size int64
	once sync.Once
	l    loader
}

func NewArchive(r iio.Reader, size int64, l loader) *archive {
	return &archive{
		r:    r,
		s:    newStorage(nil),
		size: size,
		l:    l,
	}
}

func (fs *archive) loadOnce() error {
	var errOut error
	fs.once.Do(func() {
		files, err := fs.l.getFiles(fs.r, fs.size)
		if err != nil {
			errOut = err
			return
		}

		for name, file := range files {
			if err := fs.s.Add(file, name); err != nil {
				errOut = err
				return
			}
		}
	})

	return errOut
}

func (fs *archive) Open(filename string) (File, error) {
	if filename == string(os.PathSeparator) {
		return &Dir{}, nil
	}

	if err := fs.loadOnce(); err != nil {
		return nil, err
	}

	return fs.s.Get(filename)
}

func (fs *archive) ReadDir(path string) (map[string]File, error) {
	if err := fs.loadOnce(); err != nil {
		return nil, err
	}

	return fs.s.Children(path)
}

var _ File = &ArchiveFile{}

func NewArchiveFile(readerFunc func() (iio.Reader, error), len int64) *ArchiveFile {
	return &ArchiveFile{
		readerFunc: readerFunc,
		len:        len,
	}
}

type ArchiveFile struct {
	readerFunc func() (iio.Reader, error)
	reader     iio.Reader
	len        int64
}

func (d *ArchiveFile) load() error {
	if d.reader != nil {
		return nil
	}
	r, err := d.readerFunc()
	if err != nil {
		return err
	}

	d.reader = r

	return nil
}

func (d *ArchiveFile) Size() int64 {
	return d.len
}

func (d *ArchiveFile) IsDir() bool {
	return false
}

func (d *ArchiveFile) Close() (err error) {
	if d.reader != nil {
		err = d.reader.Close()
		d.reader = nil
	}

	return
}

func (d *ArchiveFile) Read(p []byte) (n int, err error) {
	if err := d.load(); err != nil {
		return 0, err
	}

	return d.reader.Read(p)
}

func (d *ArchiveFile) ReadAt(p []byte, off int64) (n int, err error) {
	if err := d.load(); err != nil {
		return 0, err
	}

	return d.reader.ReadAt(p, off)
}
