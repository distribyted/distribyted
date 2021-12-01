package fs

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/anacrolix/missinggo/v2"
	"github.com/anacrolix/torrent"
	"github.com/distribyted/distribyted/iio"
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
				readerFunc: file.NewReader,
				len:        file.Length(),
				timeout:    fs.readTimeout,
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

type reader interface {
	iio.Reader
	missinggo.ReadContexter
}

type readAtWrapper struct {
	timeout int
	mu      sync.Mutex

	torrent.Reader
	io.ReaderAt
	io.Closer
}

func newReadAtWrapper(r torrent.Reader, timeout int) reader {
	return &readAtWrapper{Reader: r, timeout: timeout}
}

func (rw *readAtWrapper) ReadAt(p []byte, off int64) (int, error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	_, err := rw.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return readAtLeast(rw, rw.timeout, p, len(p))
}

func readAtLeast(r missinggo.ReadContexter, timeout int, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, io.ErrShortBuffer
	}
	for n < min && err == nil {
		var nn int

		ctx, cancel := context.WithCancel(context.Background())
		timer := time.AfterFunc(
			time.Duration(timeout)*time.Second,
			func() {
				cancel()
			},
		)

		nn, err = r.ReadContext(ctx, buf[n:])
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

func (rw *readAtWrapper) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.Reader.Close()
}

var _ File = &torrentFile{}

type torrentFile struct {
	readerFunc func() torrent.Reader
	reader     reader
	len        int64
	timeout    int
}

func (d *torrentFile) load() {
	if d.reader != nil {
		return
	}
	d.reader = newReadAtWrapper(d.readerFunc(), d.timeout)
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

func (d *torrentFile) ReadAt(p []byte, off int64) (n int, err error) {
	d.load()
	return d.reader.ReadAt(p, off)
}
