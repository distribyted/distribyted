package fs

import (
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/distribyted/distribyted/iio"
)

var _ Filesystem = &Torrent{}

type Torrent struct {
	mu     sync.Mutex
	ts     map[string]*torrent.Torrent
	s      *storage
	loaded bool
}

func NewTorrent() *Torrent {
	return &Torrent{
		s:  newStorage(SupportedFactories),
		ts: make(map[string]*torrent.Torrent),
	}
}

func (fs *Torrent) AddTorrent(t *torrent.Torrent) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.loaded = false
	fs.ts[t.InfoHash().HexString()] = t
}

func (fs *Torrent) RemoveTorrent(h string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.s.Clear()

	fs.loaded = false

	delete(fs.ts, h)
}

func (fs *Torrent) load() {
	if fs.loaded {
		return
	}

	fs.mu.Lock()
	defer fs.mu.Unlock()

	for _, t := range fs.ts {
		<-t.GotInfo()
		for _, file := range t.Files() {
			fs.s.Add(&torrentFile{readerFunc: file.NewReader, len: file.Length()}, file.Path())
		}
	}

	fs.loaded = true
}

func (fs *Torrent) Open(filename string) (File, error) {
	fs.load()
	return fs.s.Get(filename)
}

func (fs *Torrent) ReadDir(path string) (map[string]File, error) {
	fs.load()
	return fs.s.Children(path), nil
}

var _ File = &torrentFile{}

type torrentFile struct {
	readerFunc func() torrent.Reader
	reader     iio.Reader
	len        int64
}

func (d *torrentFile) load() {
	if d.reader != nil {
		return
	}

	d.reader = iio.NewReadAtWrapper(d.readerFunc())
}

func (d *torrentFile) Size() int64 {
	return d.len
}

func (d *torrentFile) IsDir() bool {
	return false
}

func (d *torrentFile) Close() error {
	var err error
	if d.reader != nil {
		err = d.reader.Close()
	}

	d.reader = nil

	return err
}

func (d *torrentFile) Read(p []byte) (n int, err error) {
	d.load()
	return d.reader.Read(p)
}

func (d *torrentFile) ReadAt(p []byte, off int64) (n int, err error) {
	d.load()
	return d.reader.ReadAt(p, off)
}
