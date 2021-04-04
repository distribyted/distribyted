package torrent

import (
	"fmt"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fs"
	"github.com/distribyted/distribyted/stats"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	c *torrent.Client
	s *stats.Torrent

	fssMu sync.Mutex
	fss   map[string]fs.Filesystem
}

func NewHandler(c *torrent.Client, s *stats.Torrent) *Handler {
	return &Handler{
		c:   c,
		s:   s,
		fss: make(map[string]fs.Filesystem),
	}
}

func (s *Handler) Load(route string, ts []*config.Torrent) error {
	var torrents []*torrent.Torrent
	for _, mpcTorrent := range ts {
		var t *torrent.Torrent
		var err error

		switch {
		case mpcTorrent.MagnetURI != "":
			t, err = s.c.AddMagnet(mpcTorrent.MagnetURI)
		case mpcTorrent.TorrentPath != "":
			t, err = s.c.AddTorrentFromFile(mpcTorrent.TorrentPath)
		default:
			err = fmt.Errorf("no magnet URI or torrent path provided")
		}
		if err != nil {
			return err
		}

		// only get info if name is not available
		if t.Name() == "" {
			log.Info().Str("hash", t.InfoHash().String()).Msg("getting torrent info")
			<-t.GotInfo()
		}

		s.s.Add(route, t)
		torrents = append(torrents, t)

		log.Info().Str("name", t.Name()).Str("route", route).Msg("torrent added to mountpoint")
	}

	folder := "/" + route

	s.fssMu.Lock()
	defer s.fssMu.Unlock()
	s.fss[folder] = fs.NewTorrent(torrents)

	return nil
}

func (s *Handler) Fileststems() map[string]fs.Filesystem {
	return s.fss
}

func (s *Handler) RemoveAll() error {
	s.fssMu.Lock()
	defer s.fssMu.Unlock()

	s.fss = make(map[string]fs.Filesystem)
	s.s.RemoveAll()
	return nil
}
