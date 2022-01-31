package vfs

import (
	"io"
	"io/fs"
	"os"
	"path"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var _ fs.FS = &Torrent{}

type Torrent struct {
	t *torrent.Torrent

	loadOnce sync.Once
	files    map[string]*torrent.File

	fileInfoList *FileInfoList
}

func NewTorrent(t *torrent.Torrent) *Torrent {
	return &Torrent{
		t:            t,
		fileInfoList: NewFileInfoList(),
		files:        make(map[string]*torrent.File),
	}
}

func (vfs *Torrent) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	vfs.load()

	fi := vfs.fileInfoList.Get(name)

	// special case for root
	if name == "." {
		fi = newDirFileInfo(name)
	}

	if fi == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	if fi.IsDir() {
		return newDir(fi, fileInfosToDirEntries(vfs.fileInfoList.Lookup(name))), nil
	}

	tf, ok := vfs.files[name]
	if !ok {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	tr := NewReadAtWrapper(tf.NewReader())
	rc := &ReaderCloser{
		Reader: io.NewSectionReader(tr, 0, tf.Length()),
		Closer: tr,
	}

	return newFile(fi, rc), nil
}

func (vfs *Torrent) load() {
	vfs.loadOnce.Do(func() {
		<-vfs.t.GotInfo()
		dirs := make(map[string]bool)
		for _, file := range vfs.t.Files() {
			dir, _ := split(file.Path())
			fi := newTorrentFileInfo(file)
			vfs.fileInfoList.Add(fi)

			vfs.files[file.Path()] = file
			dirs[dir] = true
		}

		for k := range dirs {
			vfs.fileInfoList.Add(newDirFileInfo(k))
		}
	})
}

var _ FileInfoPath = &torrentFileInfo{}

func newTorrentFileInfo(f *torrent.File) *torrentFileInfo {
	return &torrentFileInfo{
		f: f,
	}
}

type torrentFileInfo struct {
	f *torrent.File
}

func (fi *torrentFileInfo) Path() string {
	return fi.f.Path()
}
func (fi *torrentFileInfo) Name() string {
	base := path.Base(fi.f.Path())
	return base
}
func (fi *torrentFileInfo) Size() int64 {
	return fi.f.Length()
}
func (fi *torrentFileInfo) Mode() fs.FileMode {
	return os.ModeTemporary
}
func (fi *torrentFileInfo) ModTime() time.Time {
	return time.Time{}
}
func (fi *torrentFileInfo) IsDir() bool {
	return false
}
func (fi *torrentFileInfo) Sys() interface{} {
	return nil
}
