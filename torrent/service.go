package torrent

import (
	"errors"
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/distribyted/distribyted/fs"
	"github.com/distribyted/distribyted/torrent/loader"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Service struct {
	c *torrent.Client

	s *Stats

	mu  sync.Mutex
	fss map[string]fs.Filesystem

	cfgLoader loader.Loader
	db        loader.LoaderAdder

	log                     zerolog.Logger
	addTimeout, readTimeout int
}

func NewService(cfg loader.Loader, db loader.LoaderAdder, stats *Stats, c *torrent.Client, addTimeout, readTimeout int) *Service {
	l := log.Logger.With().Str("component", "torrent-service").Logger()
	return &Service{
		log:         l,
		s:           stats,
		c:           c,
		fss:         make(map[string]fs.Filesystem),
		cfgLoader:   cfg,
		db:          db,
		addTimeout:  addTimeout,
		readTimeout: readTimeout,
	}
}

func (s *Service) Load() (map[string]fs.Filesystem, error) {
	// Load from config
	s.log.Info().Msg("adding torrents from configuration")
	if err := s.load(s.cfgLoader); err != nil {
		return nil, err
	}

	// Load from DB
	s.log.Info().Msg("adding torrents from database")
	return s.fss, s.load(s.db)
}

func (s *Service) load(l loader.Loader) error {
	list, err := l.ListMagnets()
	if err != nil {
		return err
	}
	for r, ms := range list {
		s.addRoute(r)
		for _, m := range ms {
			if err := s.addMagnet(r, m); err != nil {
				return err
			}
		}
	}

	list, err = l.ListTorrentPaths()
	if err != nil {
		return err
	}
	for r, ms := range list {
		s.addRoute(r)
		for _, p := range ms {
			if err := s.addTorrentPath(r, p); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) AddMagnet(r, m string) error {
	if err := s.addMagnet(r, m); err != nil {
		return err
	}

	// Add to db
	return s.db.AddMagnet(r, m)
}

func (s *Service) addTorrentPath(r, p string) error {
	// Add to client
	t, err := s.c.AddTorrentFromFile(p)
	if err != nil {
		return err
	}

	return s.addTorrent(r, t)
}

func (s *Service) addMagnet(r, m string) error {
	// Add to client
	t, err := s.c.AddMagnet(m)
	if err != nil {
		return err
	}

	return s.addTorrent(r, t)

}

func (s *Service) addRoute(r string) {
	s.s.AddRoute(r)

	// Add to filesystems
	folder := path.Join("/", r)
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.fss[folder]
	if !ok {
		s.fss[folder] = fs.NewTorrent(s.readTimeout)
	}
}

func (s *Service) addTorrent(r string, t *torrent.Torrent) error {
	// only get info if name is not available
	if t.Info() == nil {
		s.log.Info().Str("hash", t.InfoHash().String()).Msg("getting torrent info")
		select {
		case <-time.After(time.Duration(s.addTimeout) * time.Second):
			s.log.Error().Str("hash", t.InfoHash().String()).Msg("timeout getting torrent info")
			return errors.New("timeout getting torrent info")
		case <-t.GotInfo():
			s.log.Info().Str("hash", t.InfoHash().String()).Msg("obtained torrent info")
		}

	}

	// Add to stats
	s.s.Add(r, t)

	// Add to filesystems
	folder := path.Join("/", r)
	s.mu.Lock()
	defer s.mu.Unlock()

	tfs, ok := s.fss[folder].(*fs.Torrent)
	if !ok {
		return errors.New("error adding torrent to filesystem")
	}

	tfs.AddTorrent(t)
	s.log.Info().Str("name", t.Info().Name).Str("route", r).Msg("torrent added")

	return nil
}

func (s *Service) RemoveFromHash(r, h string) error {
	// Remove from db
	deleted, err := s.db.RemoveFromHash(r, h)
	if err != nil {
		return err
	}

	if !deleted {
		return fmt.Errorf("element with hash %v on route %v cannot be removed", h, r)
	}

	// Remove from stats
	s.s.Del(r, h)

	// Remove from fs
	folder := path.Join("/", r)

	tfs, ok := s.fss[folder].(*fs.Torrent)
	if !ok {
		return errors.New("error removing torrent from filesystem")
	}

	tfs.RemoveTorrent(h)

	// Remove from client
	var mh metainfo.Hash
	if err := mh.FromHexString(h); err != nil {
		return err
	}

	t, ok := s.c.Torrent(metainfo.NewHashFromHex(h))
	if ok {
		t.Drop()
	}

	return nil
}
