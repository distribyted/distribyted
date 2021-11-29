package torrent

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

var ErrTorrentNotFound = errors.New("torrent not found")

type PieceStatus string

const (
	Checking PieceStatus = "H"
	Partial  PieceStatus = "P"
	Complete PieceStatus = "C"
	Waiting  PieceStatus = "W"
	Error    PieceStatus = "?"
)

type PieceChunk struct {
	Status    PieceStatus `json:"status"`
	NumPieces int         `json:"numPieces"`
}

type TorrentStats struct {
	Name            string        `json:"name"`
	Hash            string        `json:"hash"`
	DownloadedBytes int64         `json:"downloadedBytes"`
	UploadedBytes   int64         `json:"uploadedBytes"`
	Peers           int           `json:"peers"`
	Seeders         int           `json:"seeders"`
	TimePassed      float64       `json:"timePassed"`
	PieceChunks     []*PieceChunk `json:"pieceChunks"`
	TotalPieces     int           `json:"totalPieces"`
	PieceSize       int64         `json:"pieceSize"`
}

type byName []*TorrentStats

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type GlobalTorrentStats struct {
	DownloadedBytes int64   `json:"downloadedBytes"`
	UploadedBytes   int64   `json:"uploadedBytes"`
	TimePassed      float64 `json:"timePassed"`
}

type RouteStats struct {
	Name         string          `json:"name"`
	TorrentStats []*TorrentStats `json:"torrentStats"`
}

type ByName []*RouteStats

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type stat struct {
	totalDownloadBytes int64
	downloadBytes      int64
	totalUploadBytes   int64
	uploadBytes        int64
	peers              int
	seeders            int
	time               time.Time
}

type Stats struct {
	mut             sync.Mutex
	torrents        map[string]*torrent.Torrent
	torrentsByRoute map[string]map[string]*torrent.Torrent
	previousStats   map[string]*stat

	gTime time.Time
}

func NewStats() *Stats {
	return &Stats{
		gTime:           time.Now(),
		torrents:        make(map[string]*torrent.Torrent),
		torrentsByRoute: make(map[string]map[string]*torrent.Torrent),
		previousStats:   make(map[string]*stat),
	}
}

func (s *Stats) AddRoute(route string) {
	_, ok := s.torrentsByRoute[route]
	if !ok {
		s.torrentsByRoute[route] = make(map[string]*torrent.Torrent)
	}
}

func (s *Stats) Add(route string, t *torrent.Torrent) {
	s.mut.Lock()
	defer s.mut.Unlock()

	h := t.InfoHash().String()

	s.torrents[h] = t
	s.previousStats[h] = &stat{}

	_, ok := s.torrentsByRoute[route]
	if !ok {
		s.torrentsByRoute[route] = make(map[string]*torrent.Torrent)
	}

	s.torrentsByRoute[route][h] = t
}

func (s *Stats) Del(route, hash string) {
	s.mut.Lock()
	defer s.mut.Unlock()
	delete(s.torrents, hash)
	delete(s.previousStats, hash)
	ts, ok := s.torrentsByRoute[route]
	if !ok {
		return
	}

	delete(ts, hash)
}

func (s *Stats) Stats(hash string) (*TorrentStats, error) {
	s.mut.Lock()
	defer s.mut.Unlock()

	t, ok := s.torrents[hash]
	if !(ok) {
		return nil, ErrTorrentNotFound
	}

	now := time.Now()

	return s.stats(now, t, true), nil
}

func (s *Stats) RoutesStats() []*RouteStats {
	s.mut.Lock()
	defer s.mut.Unlock()

	now := time.Now()

	var out []*RouteStats
	for r, tl := range s.torrentsByRoute {
		var tStats []*TorrentStats
		for _, t := range tl {
			ts := s.stats(now, t, true)
			tStats = append(tStats, ts)
		}

		sort.Sort(byName(tStats))

		rs := &RouteStats{
			Name:         r,
			TorrentStats: tStats,
		}
		out = append(out, rs)
	}

	return out
}

func (s *Stats) GlobalStats() *GlobalTorrentStats {
	s.mut.Lock()
	defer s.mut.Unlock()

	now := time.Now()

	var totalDownload int64
	var totalUpload int64
	for _, torrent := range s.torrents {
		tStats := s.stats(now, torrent, false)
		totalDownload += tStats.DownloadedBytes
		totalUpload += tStats.UploadedBytes
	}

	timePassed := now.Sub(s.gTime)
	s.gTime = now

	return &GlobalTorrentStats{
		DownloadedBytes: totalDownload,
		UploadedBytes:   totalUpload,
		TimePassed:      timePassed.Seconds(),
	}
}

func (s *Stats) stats(now time.Time, t *torrent.Torrent, chunks bool) *TorrentStats {
	ts := &TorrentStats{}
	prev, ok := s.previousStats[t.InfoHash().String()]
	if !ok {
		return &TorrentStats{}
	}
	if s.returnPreviousMeasurements(now) {
		ts.DownloadedBytes = prev.downloadBytes
		ts.UploadedBytes = prev.uploadBytes
	} else {
		st := t.Stats()
		rd := st.BytesReadData.Int64()
		wd := st.BytesWrittenData.Int64()
		ist := &stat{
			downloadBytes:      rd - prev.totalDownloadBytes,
			uploadBytes:        wd - prev.totalUploadBytes,
			totalDownloadBytes: rd,
			totalUploadBytes:   wd,
			time:               now,
			peers:              st.TotalPeers,
			seeders:            st.ConnectedSeeders,
		}

		ts.DownloadedBytes = ist.downloadBytes
		ts.UploadedBytes = ist.uploadBytes
		ts.Peers = ist.peers
		ts.Seeders = ist.seeders

		s.previousStats[t.InfoHash().String()] = ist
	}

	ts.TimePassed = now.Sub(prev.time).Seconds()
	var totalPieces int
	if chunks {
		var pch []*PieceChunk
		for _, psr := range t.PieceStateRuns() {
			var s PieceStatus
			switch {
			case psr.Checking:
				s = Checking
			case psr.Partial:
				s = Partial
			case psr.Complete:
				s = Complete
			case !psr.Ok:
				s = Error
			default:
				s = Waiting
			}

			pch = append(pch, &PieceChunk{
				Status:    s,
				NumPieces: psr.Length,
			})
			totalPieces += psr.Length
		}
		ts.PieceChunks = pch
	}

	ts.Hash = t.InfoHash().String()
	ts.Name = t.Name()
	ts.TotalPieces = totalPieces

	if t.Info() != nil {
		ts.PieceSize = t.Info().PieceLength
	}

	return ts
}

const gap time.Duration = 2 * time.Second

func (s *Stats) returnPreviousMeasurements(now time.Time) bool {
	return now.Sub(s.gTime) < gap
}
