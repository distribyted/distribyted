package torrent

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/storage"
	"github.com/distribyted/distribyted/config"
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ServerState int

const (
	UNKNOWN ServerState = iota
	SEEDING
	READING
	UPDATING
	STOPPED
	ERROR
)

func (ss ServerState) String() string {
	return [...]string{"Unknown", "Seeding", "Reading", "Updating", "Stopped", "Error"}[ss]
}

type ServerInfo struct {
	Magnet    string `json:"magnetUri"`
	UpdatedAt int64  `json:"updatedAt"`
	Name      string `json:"name"`
	Folder    string `json:"folder"`
	State     string `json:"state"`
	Peers     int    `json:"peers"`
	Seeds     int    `json:"seeds"`
}

type Server struct {
	cfg *config.Server
	log zerolog.Logger

	fw *fsnotify.Watcher

	eventsCount uint64

	c  *torrent.Client
	pc storage.PieceCompletion

	mu sync.RWMutex
	t  *torrent.Torrent
	si ServerInfo
}

func NewServer(c *torrent.Client, pc storage.PieceCompletion, cfg *config.Server) *Server {
	l := log.Logger.With().Str("component", "server").Str("name", cfg.Name).Logger()

	return &Server{
		cfg: cfg,
		log: l,
		c:   c,
		pc:  pc,
	}
}

func (s *Server) Start() error {
	s.log.Info().Msg("starting new server folder")
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(s.cfg.Path, 0744); err != nil {
		return fmt.Errorf("error creating server folder: %s. Error: %w", s.cfg.Path, err)
	}

	if err := filepath.Walk(s.cfg.Path,
		func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsDir() {
				s.log.Debug().Str("folder", path).Msg("adding new folder")
				return w.Add(path)
			}

			return nil
		}); err != nil {
		return err
	}

	s.fw = w
	go func() {
		if err := s.makeMagnet(); err != nil {
			s.updateState(ERROR)
			s.log.Error().Err(err).Msg("error generating magnet on start")
		}

		s.watch()
	}()

	go func() {
		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}

				s.log.Info().Str("file", event.Name).Str("op", event.Op.String()).Msg("file changed inside server folder")
				atomic.AddUint64(&s.eventsCount, 1)
			case err, ok := <-w.Errors:
				if !ok {
					return
				}

				s.updateState(STOPPED)
				s.log.Error().Err(err).Msg("error watching server folder")
			}
		}
	}()

	s.log.Info().Msg("server folder started")

	return nil
}

func (s *Server) watch() {
	s.log.Info().Msg("starting watcher")
	for range time.Tick(time.Second * 5) {
		if s.eventsCount == 0 {
			continue
		}

		ec := s.eventsCount
		if err := s.makeMagnet(); err != nil {
			s.updateState(ERROR)
			s.log.Error().Err(err).Msg("error generating magnet")
		}

		atomic.AddUint64(&s.eventsCount, -ec)
	}
}

func (s *Server) makeMagnet() error {

	s.log.Info().Msg("starting serving new torrent")

	info := metainfo.Info{
		PieceLength: 2 << 18,
	}

	s.updateState(READING)

	if err := info.BuildFromFilePath(s.cfg.Path); err != nil {
		return err
	}

	s.updateState(UPDATING)

	if len(info.Files) == 0 {
		s.mu.Lock()
		s.si.Magnet = ""
		s.si.Folder = s.cfg.Path
		s.si.Name = s.cfg.Name
		s.si.UpdatedAt = time.Now().Unix()
		s.mu.Unlock()
		s.log.Info().Msg("not creating magnet from empty folder")

		s.updateState(STOPPED)
		return nil
	}

	mi := metainfo.MetaInfo{
		InfoBytes: bencode.MustMarshal(info),
	}

	ih := mi.HashInfoBytes()

	to, _ := s.c.AddTorrentOpt(torrent.AddTorrentOpts{
		InfoHash: ih,
		Storage: storage.NewFileOpts(storage.NewFileClientOpts{
			ClientBaseDir: s.cfg.Path,
			FilePathMaker: func(opts storage.FilePathMakerOpts) string {
				return filepath.Join(opts.File.Path...)
			},
			TorrentDirMaker: nil,
			PieceCompletion: s.pc,
		}),
	})

	tks := s.trackers()

	err := to.MergeSpec(&torrent.TorrentSpec{
		InfoBytes: mi.InfoBytes,

		Trackers: [][]string{tks},
	})
	if err != nil {
		return err
	}

	m := metainfo.Magnet{
		InfoHash:    ih,
		DisplayName: s.cfg.Name,
		Trackers:    tks,
	}

	s.mu.Lock()
	s.t = to
	s.si.Magnet = m.String()
	s.si.Folder = s.cfg.Path
	s.si.Name = s.cfg.Name
	s.si.UpdatedAt = time.Now().Unix()
	s.mu.Unlock()
	s.updateState(SEEDING)

	s.log.Info().Str("hash", ih.HexString()).Msg("new torrent is ready")

	return nil
}

func (s *Server) updateState(ss ServerState) {
	s.mu.Lock()
	s.si.State = ss.String()
	s.mu.Unlock()
}

func (s *Server) trackers() []string {
	// TODO load trackers from URL too
	return s.cfg.Trackers
}

func (s *Server) Close() error {
	if s.fw == nil {
		return nil
	}
	return s.fw.Close()
}

func (s *Server) Info() *ServerInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.t != nil {
		st := s.t.Stats()
		s.si.Peers = st.TotalPeers
		s.si.Seeds = st.ConnectedSeeders
	}

	return &s.si
}
