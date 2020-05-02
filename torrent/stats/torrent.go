package stats

import (
	"errors"
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

var ErrTorrentNotFound = errors.New("torrent not found")

type PieceStatus rune

const (
	Checking PieceStatus = 'H'
	Partial  PieceStatus = 'P'
	Complete PieceStatus = 'C'
	Waiting  PieceStatus = 'W'
	Error    PieceStatus = '?'
)

type PieceChunk struct {
	Status    PieceStatus
	NumPieces int
}

type TorrentStats struct {
	DownloadedBytes int64
	UploadedBytes   int64
	TimePassed      time.Duration
	PieceChunks     []*PieceChunk
}

type GlobalTorrentStats struct {
	DownloadedBytes int64
	UploadedBytes   int64
	TimePassed      time.Duration
}

func (s *GlobalTorrentStats) speed(bytes int64) float64 {
	var bs float64
	t := s.TimePassed.Seconds()
	if t != 0 {
		bs = float64(bytes) / t
	}

	return bs
}

func (s *GlobalTorrentStats) DownloadSpeed() string {
	return fmt.Sprintf(" %s/s", humanize.IBytes(uint64(s.speed(s.DownloadedBytes))))
}

func (s *GlobalTorrentStats) UploadSpeed() string {
	return fmt.Sprintf(" %s/s", humanize.IBytes(uint64(s.speed(s.UploadedBytes))))
}

type stats struct {
	totalDownloadBytes int64
	downloadBytes      int64
	totalUploadBytes   int64
	uploadBytes        int64
	time               time.Time
}

type Torrent struct {
	torrents      map[string]*torrent.Torrent
	previousStats map[string]*stats

	gTime time.Time
}

func NewTorrent() *Torrent {
	return &Torrent{
		gTime:         time.Now(),
		torrents:      make(map[string]*torrent.Torrent),
		previousStats: make(map[string]*stats),
	}
}

func (s *Torrent) Add(t *torrent.Torrent) {
	s.torrents[t.InfoHash().String()] = t
	s.previousStats[t.InfoHash().String()] = &stats{}
}

func (s *Torrent) Torrent(hash string) (*TorrentStats, error) {
	t, ok := s.torrents[hash]
	if !(ok) {
		return nil, ErrTorrentNotFound
	}

	now := time.Now()

	return s.stats(now, t, true), nil
}

func (s *Torrent) List() []string {
	var result []string
	for hash := range s.torrents {
		result = append(result, hash)
	}

	return result
}

func (s *Torrent) Global() *GlobalTorrentStats {
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
		TimePassed:      timePassed,
	}
}

func (s *Torrent) stats(now time.Time, t *torrent.Torrent, chunks bool) *TorrentStats {
	ts := &TorrentStats{}
	prev := s.previousStats[t.InfoHash().String()]
	if s.returnPreviousMeasurements(now) {
		ts.DownloadedBytes = prev.downloadBytes
		ts.UploadedBytes = prev.uploadBytes
	} else {
		st := t.Stats()

		rd := st.BytesReadData.Int64()
		wd := st.BytesWrittenData.Int64()
		ist := &stats{
			downloadBytes:      rd - prev.totalDownloadBytes,
			uploadBytes:        wd - prev.totalUploadBytes,
			totalDownloadBytes: rd,
			totalUploadBytes:   wd,
			time:               now,
		}

		ts.DownloadedBytes = ist.downloadBytes
		ts.UploadedBytes = ist.uploadBytes

		s.previousStats[t.InfoHash().String()] = ist
	}

	ts.TimePassed = now.Sub(prev.time)

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
		}
		ts.PieceChunks = pch
	}

	return ts
}

const gap time.Duration = 2 * time.Second

func (s *Torrent) returnPreviousMeasurements(now time.Time) bool {
	return now.Sub(s.gTime) < gap
}
