package vfs

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"sync"
	"time"
)

func newFile(fi fs.FileInfo, r io.ReadCloser) *File {
	return &File{
		fi:         fi,
		ReadCloser: r,
	}
}

func newDir(fi fs.FileInfo, dirs []fs.DirEntry) *File {
	return &File{
		fi:         fi,
		ReadCloser: ioutil.NopCloser(nil),
		dirs:       dirs,
	}
}

var _ fs.DirEntry = &DirEntry{}

type DirEntry struct {
	fi fs.FileInfo
}

func newDirEntry(fi fs.FileInfo) *DirEntry {
	return &DirEntry{
		fi: fi,
	}
}

func (de *DirEntry) Name() string {
	return de.fi.Name()
}

func (de *DirEntry) IsDir() bool {
	return de.fi.IsDir()
}

func (de *DirEntry) Type() fs.FileMode {
	return de.fi.Mode().Type()
}

func (de *DirEntry) Info() (fs.FileInfo, error) {
	return de.fi, nil
}

var _ fs.File = &File{}
var _ fs.ReadDirFile = &File{}

type File struct {
	fi fs.FileInfo
	io.ReadCloser

	offset int
	dirs   []fs.DirEntry
}

func (f *File) Stat() (fs.FileInfo, error) {
	return f.fi, nil
}

func (f *File) ReadDir(count int) ([]fs.DirEntry, error) {
	if !f.fi.IsDir() {
		return nil, errors.New("not a directory")
	}

	n := len(f.dirs) - f.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 {
		if count <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}

	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = f.dirs[f.offset+i]
	}
	f.offset += n
	return list, nil
}

type FileInfoPath interface {
	fs.FileInfo
	Path() string
}

func fileInfosToDirEntries(fis []FileInfoPath) []os.DirEntry {
	var out []os.DirEntry
	for _, fi := range fis {
		out = append(out, newDirEntry(fi))
	}

	return out
}

type FileInfoList struct {
	fil    []FileInfoPath
	sorted bool
}

func NewFileInfoList() *FileInfoList {
	return &FileInfoList{}
}

func (fil *FileInfoList) Add(fi FileInfoPath) {
	fil.fil = append(fil.fil, fi)
	fil.sorted = false
}

func (fil *FileInfoList) Sort() {
	sort.Slice(fil.fil, func(i, j int) bool { return fileEntryLess(fil.fil[i].Path(), fil.fil[j].Path()) })
	fil.sorted = true
}

func (fil *FileInfoList) Get(path string) FileInfoPath {
	for _, fi := range fil.fil {
		if fi.Path() == path {
			return fi
		}
	}

	return nil
}

func (fil *FileInfoList) Lookup(dir string) []FileInfoPath {
	if !fil.sorted {
		fil.Sort()
	}

	fis := fil.fil
	i := sort.Search(len(fis), func(i int) bool {
		idir, _ := split(fis[i].Path())
		return idir >= dir
	})
	j := sort.Search(len(fis), func(j int) bool {
		jdir, _ := split(fis[j].Path())
		return jdir > dir
	})

	return fis[i:j]
}

func fileEntryLess(x, y string) bool {
	xdir, xelem := split(x)
	ydir, yelem := split(y)
	return xdir < ydir || xdir == ydir && xelem < yelem
}

func split(name string) (dir, elem string) {
	i := len(name) - 1
	for i >= 0 && name[i] != '/' {
		i--
	}
	if i < 0 {
		return ".", name
	}
	return name[:i], name[i+1:]
}

var _ FileInfoPath = &dirFileInfo{}

func newDirFileInfo(name string) *dirFileInfo {
	return &dirFileInfo{name: name}
}

type dirFileInfo struct {
	name string
}

func (fi *dirFileInfo) Path() string {
	return fi.name
}

func (fi *dirFileInfo) Name() string {
	return path.Base(fi.name)
}
func (fi *dirFileInfo) Size() int64 {
	return 0
}
func (fi *dirFileInfo) Mode() fs.FileMode {
	return fs.ModeDir | 0555
}
func (fi *dirFileInfo) ModTime() time.Time {
	return time.Time{}
}
func (fi *dirFileInfo) IsDir() bool {
	return true
}
func (fi *dirFileInfo) Sys() interface{} {
	return nil
}

func isSymlink(info fs.FileInfo) bool {
	return info.Mode()&os.ModeSymlink != 0
}

type ReaderCloser struct {
	io.Reader
	io.Closer
}

var _ io.ReaderAt = &ReadAtWrapper{}
var _ io.Closer = &ReadAtWrapper{}

type ReadAtWrapper struct {
	mu sync.Mutex
	io.ReadSeekCloser
}

func NewReadAtWrapper(r io.ReadSeekCloser) *ReadAtWrapper {
	return &ReadAtWrapper{ReadSeekCloser: r}
}

func (rw *ReadAtWrapper) ReadAt(p []byte, off int64) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// Get actual position
	pos, err := rw.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	defer rw.Seek(pos, io.SeekStart)

	_, err = rw.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return readAtLeast(rw, p, len(p))
}

func readAtLeast(r io.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int

		nn, err = r.Read(buf[n:])
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}
