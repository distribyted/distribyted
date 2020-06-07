package mount

import (
	"fmt"
	"log"
	"os"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/node"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type Handler struct {
	c    *torrent.Client
	s    *stats.Torrent
	opts *fs.Options

	servers map[string]*fuse.Server
}

func NewHandler(c *torrent.Client, s *stats.Torrent) *Handler {
	return &Handler{
		c:       c,
		s:       s,
		opts:    &fs.Options{},
		servers: make(map[string]*fuse.Server),
	}
}

func (s *Handler) Mount(mpc *config.MountPoint) error {
	var torrents []*torrent.Torrent
	for _, mpcTorrent := range mpc.Torrents {
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
			log.Println("getting torrent info", t.InfoHash())
			<-t.GotInfo()
		}

		s.s.Add(mpc.Path, t)
		log.Println("torrent added", t.Name())

		torrents = append(torrents, t)
	}

	// TODO change permissions
	if err := os.MkdirAll(mpc.Path, 0770); err != nil && !os.IsExist(err) {
		log.Println("UFFF", err)
		return err
	}

	node := node.NewRoot(torrents)
	server, err := fs.Mount(mpc.Path, node, s.opts)
	if err != nil {
		return err
	}

	s.servers[mpc.Path] = server

	return nil
}

func (s *Handler) Close() {
	for path, server := range s.servers {
		log.Println("unmounting", path)
		err := server.Unmount()
		if err != nil {
			log.Println("unmount failed", path, err)
		}
	}
}
