package fs

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var _ Filesystem = &Torrent{}

type Torrent struct {
	mu          sync.RWMutex
	ts          map[string]*torrent.Torrent
	s           *storage
	loaded      bool
	readTimeout int
}

func NewTorrent(readTimeout int) *Torrent {
	return &Torrent{
		s:           newStorage(SupportedFactories),
		ts:          make(map[string]*torrent.Torrent),
		readTimeout: readTimeout,
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
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	for _, t := range fs.ts {
		<-t.GotInfo()
		for _, file := range t.Files() {
			fs.s.Add(&torrentFile{
				reader:  file.NewReader(),
				len:     file.Length(),
				timeout: fs.readTimeout,
			}, file.Path())
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
	return fs.s.Children(path)
}

var _ File = &torrentFile{}

type torrentFile struct {
	mu sync.Mutex

	reader torrent.Reader
	len    int64

	timeout int
}

func (d *torrentFile) Size() int64 {
	return d.len
}

func (d *torrentFile) IsDir() bool {
	return false
}

func (d *torrentFile) Close() error {
	return d.reader.Close()
}

func (d *torrentFile) Read(p []byte) (n int, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := time.AfterFunc(
		time.Duration(d.timeout)*time.Second,
		func() {
			cancel()
		},
	)

	defer timer.Stop()

	return d.reader.ReadContext(ctx, p)
}

func (d *torrentFile) ReadAt(p []byte, off int64) (int, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.reader.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	i, err := d.readAtLeast(p, len(p))
	return i, err
}

func (d *torrentFile) readAtLeast(buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int

		ctx, cancel := context.WithCancel(context.Background())
		timer := time.AfterFunc(
			time.Duration(d.timeout)*time.Second,
			func() {
				cancel()
			},
		)

		nn, err = d.reader.ReadContext(ctx, buf[n:])
		n += nn

		timer.Stop()
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}
