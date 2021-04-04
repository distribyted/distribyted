package torrent

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/distribyted/distribyted/config"
)

func NewClient(st storage.ClientImpl, cfg *config.TorrentGlobal) (*torrent.Client, error) {
	// TODO download and upload limits
	torrentCfg := torrent.NewDefaultClientConfig()
	torrentCfg.Logger = log.Discard
	torrentCfg.Seed = true
	torrentCfg.DefaultStorage = st

	torrentCfg.DisableIPv6 = cfg.DisableIPv6

	return torrent.NewClient(torrentCfg)
}
