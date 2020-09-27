package fs

import (
	"archive/zip"
	"os"

	"github.com/ajnavarro/distribyted/iio"
)

var _ Filesystem = &Zip{}

type Zip struct {
	r    iio.Reader
	s    *storage
	size int64

	loaded bool
}

func NewZip(r iio.Reader, size int64) *Zip {
	return &Zip{
		r:    r,
		size: size,
		s:    newStorage(nil),
	}
}

func (fs *Zip) load() error {
	if fs.loaded {
		return nil
	}

	zr, err := zip.NewReader(fs.r, fs.size)
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		f := f
		if f.FileInfo().IsDir() {
			continue
		}

		err := fs.s.Add(newZipFile(
			func() (iio.Reader, error) {
				zr, err := f.Open()
				if err != nil {
					return nil, err
				}

				return iio.NewDiskTeeReader(zr)
			},
			f.FileInfo().Size(),
		), string(os.PathSeparator)+f.Name)
		if err != nil {
			return err
		}
	}

	fs.loaded = true

	return nil
}

func (fs *Zip) Open(filename string) (File, error) {
	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs.s.Get(filename)
}

func (fs *Zip) ReadDir(path string) (map[string]File, error) {
	if err := fs.load(); err != nil {
		return nil, err
	}

	return fs.s.Children(path), nil
}

var _ File = &zipFile{}

func newZipFile(readerFunc func() (iio.Reader, error), len int64) *zipFile {
	return &zipFile{
		readerFunc: readerFunc,
		len:        len,
	}
}

type zipFile struct {
	readerFunc func() (iio.Reader, error)
	reader     iio.Reader
	len        int64
}

func (d *zipFile) load() error {
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

func (d *zipFile) Size() int64 {
	return d.len
}

func (d *zipFile) IsDir() bool {
	return false
}

func (d *zipFile) Close() (err error) {
	if d.reader != nil {
		err = d.reader.Close()
		d.reader = nil
	}

	return
}

func (d *zipFile) Read(p []byte) (n int, err error) {
	if err := d.load(); err != nil {
		return 0, err
	}

	return d.reader.Read(p)
}

func (d *zipFile) ReadAt(p []byte, off int64) (n int, err error) {
	if err := d.load(); err != nil {
		return 0, err
	}

	return d.reader.ReadAt(p, off)
}
