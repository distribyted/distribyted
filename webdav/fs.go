package webdav

import (
	"context"
	"io"
	"os"
	"path"
	"sync"
	"time"

	"github.com/distribyted/distribyted/fs"
	"github.com/distribyted/distribyted/iio"
	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = &WebDAV{}

type WebDAV struct {
	fss []fs.Filesystem
}

func newFS(mFss map[string][]fs.Filesystem) *WebDAV {
	for _, fss := range mFss {
		return &WebDAV{fss: fss}
	}

	return nil
}

func (wd *WebDAV) OpenFile(ctx context.Context, name string, flag int, perm os.FileMode) (webdav.File, error) {
	p := "/" + name
	// TODO handle flag and permissions
	f, err := wd.lookupFile(p)
	if err != nil {
		return nil, err
	}

	var dirContent []os.FileInfo
	if f.IsDir() {
		dir, err := wd.listDir(p)
		if err != nil {
			return nil, err
		}

		dirContent = dir
	}

	wdf := newFile(path.Base(p), f, dirContent)
	return wdf, nil
}

func (wd *WebDAV) Stat(ctx context.Context, name string) (os.FileInfo, error) {
	p := "/" + name
	f, err := wd.lookupFile(p)
	if err != nil {
		return nil, err
	}
	fi := newFileInfo(name, f.Size(), f.IsDir())
	return fi, nil
}

func (wd *WebDAV) Mkdir(ctx context.Context, name string, perm os.FileMode) error {
	return webdav.ErrNotImplemented
}

func (wd *WebDAV) RemoveAll(ctx context.Context, name string) error {
	return webdav.ErrNotImplemented
}

func (wd *WebDAV) Rename(ctx context.Context, oldName, newName string) error {
	return webdav.ErrNotImplemented
}

func (wd *WebDAV) lookupFile(path string) (fs.File, error) {
	for _, f := range wd.fss {
		file, err := f.Open(path)
		if err == os.ErrNotExist {
			continue
		}
		if err != nil {
			return nil, err
		}

		if file != nil {
			return file, nil
		}
	}

	return nil, os.ErrNotExist
}

func (wd *WebDAV) listDir(path string) ([]os.FileInfo, error) {
	var out []os.FileInfo
	for _, ifs := range wd.fss {
		files, err := ifs.ReadDir(path)
		if err != nil {
			return nil, err
		}

		for n, f := range files {
			out = append(out, newFileInfo(n, f.Size(), f.IsDir()))
		}
	}

	return out, nil
}

var _ webdav.File = &webDAVFile{}

type webDAVFile struct {
	iio.Reader

	fi os.FileInfo

	mu sync.Mutex
	// dirPos and pos are protected by mu.
	dirPos     int
	pos        int64
	dirContent []os.FileInfo
}

func newFile(name string, f fs.File, dir []os.FileInfo) *webDAVFile {
	return &webDAVFile{
		fi:         newFileInfo(name, f.Size(), f.IsDir()),
		dirContent: dir,
		Reader:     f,
	}
}

func (wdf *webDAVFile) Readdir(count int) ([]os.FileInfo, error) {
	wdf.mu.Lock()
	defer wdf.mu.Unlock()

	if !wdf.fi.IsDir() {
		return nil, os.ErrInvalid
	}

	old := wdf.dirPos
	if old >= len(wdf.dirContent) {
		// The os.File Readdir docs say that at the end of a directory,
		// the error is io.EOF if count > 0 and nil if count <= 0.
		if count > 0 {
			return nil, io.EOF
		}
		return nil, nil
	}
	if count > 0 {
		wdf.dirPos += count
		if wdf.dirPos > len(wdf.dirContent) {
			wdf.dirPos = len(wdf.dirContent)
		}
	} else {
		wdf.dirPos = len(wdf.dirContent)
		old = 0
	}

	return wdf.dirContent[old:wdf.dirPos], nil
}

func (wdf *webDAVFile) Stat() (os.FileInfo, error) {
	return wdf.fi, nil
}

func (wdf *webDAVFile) Read(p []byte) (int, error) {
	wdf.mu.Lock()
	defer wdf.mu.Unlock()

	n, err := wdf.Reader.ReadAt(p, wdf.pos)
	wdf.pos += int64(n)

	return n, err
}

func (wdf *webDAVFile) Seek(offset int64, whence int) (int64, error) {
	wdf.mu.Lock()
	defer wdf.mu.Unlock()

	switch whence {
	case io.SeekStart:
		wdf.pos = offset
		break
	case io.SeekCurrent:
		wdf.pos = wdf.pos + offset
		break
	case io.SeekEnd:
		wdf.pos = wdf.fi.Size() + offset
		break
	}

	return wdf.pos, nil
}

func (wdf *webDAVFile) Write(p []byte) (n int, err error) {
	return 0, webdav.ErrNotImplemented
}

type webDAVFileInfo struct {
	name  string
	size  int64
	isDir bool
}

func newFileInfo(name string, size int64, isDir bool) *webDAVFileInfo {
	return &webDAVFileInfo{
		name:  name,
		size:  size,
		isDir: isDir,
	}
}

func (wdfi *webDAVFileInfo) Name() string {
	return wdfi.name
}

func (wdfi *webDAVFileInfo) Size() int64 {
	return wdfi.size
}

func (wdfi *webDAVFileInfo) Mode() os.FileMode {
	if wdfi.isDir {
		return 0555 | os.ModeDir
	}

	return 0555
}

func (wdfi *webDAVFileInfo) ModTime() time.Time {
	// TODO fix it
	return time.Now()
}

func (wdfi *webDAVFileInfo) IsDir() bool {
	return wdfi.isDir
}

func (wdfi *webDAVFileInfo) Sys() interface{} {
	return nil
}
