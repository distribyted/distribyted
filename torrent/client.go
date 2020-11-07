package torrent

import (
	"github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
)

func NewClient(st storage.ClientImpl) (*torrent.Client, error) {
	// TODO download and upload limits
	torrentCfg := torrent.NewDefaultClientConfig()
	torrentCfg.Logger = log.Discard
	torrentCfg.Seed = true
	torrentCfg.DisableTCP = true
	torrentCfg.DefaultStorage = st

	return torrent.NewClient(torrentCfg)

}
