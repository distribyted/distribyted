package torrent

import (
	"fmt"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fs"
	"github.com/distribyted/distribyted/stats"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	c *torrent.Client
	s *stats.Torrent

	fssMu sync.Mutex
	fss   map[string][]fs.Filesystem
}

func NewHandler(c *torrent.Client, s *stats.Torrent) *Handler {
	return &Handler{
		c:   c,
		s:   s,
		fss: make(map[string][]fs.Filesystem),
	}
}

func (s *Handler) Load(path string, ts []*config.Torrent) error {
	var torrents []fs.Filesystem
	for _, mpcTorrent := range ts {
		var t *torrent.Torrent
		var err error

		switch {
		case mpcTorrent.MagnetURI != "":
			t, err = s.c.AddMagnet(mpcTorrent.MagnetURI)
			break
		case mpcTorrent.TorrentPath != "":
			t, err = s.c.AddTorrentFromFile(mpcTorrent.TorrentPath)
			break
		default:
			err = fmt.Errorf("no magnet URI or torrent path provided")
		}
		if err != nil {
			return err
		}

		// only get info if name is not available
		if t.Name() == "" {
			log.WithField("hash", t.InfoHash()).Info("getting torrent info")
			<-t.GotInfo()
		}

		s.s.Add(path, t)
		torrents = append(torrents, fs.NewTorrent(t))

		log.WithField("name", t.Name()).WithField("path", path).Info("torrent added to mountpoint")
	}

	folder := path

	s.fssMu.Lock()
	defer s.fssMu.Unlock()
	s.fss[folder] = torrents

	return nil
}

func (s *Handler) Fileststems() map[string][]fs.Filesystem {
	return s.fss
}

func (s *Handler) RemoveAll() error {
	s.fssMu.Lock()
	defer s.fssMu.Unlock()

	s.fss = make(map[string][]fs.Filesystem)
	s.s.RemoveAll()
	return nil
}
