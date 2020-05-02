package mount

import (
	"log"
	"os"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/node"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/panjf2000/ants/v2"
)

type Torrent struct {
	c    *torrent.Client
	s    *stats.Torrent
	opts *fs.Options

	pool    *ants.Pool
	servers map[string]*fuse.Server
}

func NewTorrent(c *torrent.Client, pool *ants.Pool, s *stats.Torrent) *Torrent {
	return &Torrent{
		c:       c,
		s:       s,
		opts:    &fs.Options{},
		pool:    pool,
		servers: make(map[string]*fuse.Server),
	}
}

func (s *Torrent) Mount(mpc *config.MountPoint) error {
	var torrents []*torrent.Torrent
	for _, magnet := range mpc.Magnets {
		t, err := s.c.AddMagnet(magnet.URI)
		if err != nil {
			return err
		}
		log.Println("getting torrent info", t.Name())
		<-t.GotInfo()

		s.s.Add(t)

		log.Println("torrent info obtained", t.Name())
		torrents = append(torrents, t)
	}

	// TODO change permissions
	if err := os.MkdirAll(mpc.Path, 0770); err != nil && !os.IsExist(err) {
		log.Println("UFFF", err)
		return err
	}

	node := node.NewRoot(torrents, s.pool)
	server, err := fs.Mount(mpc.Path, node, s.opts)
	if err != nil {
		return err
	}

	s.servers[mpc.Path] = server

	return nil
}

func (s *Torrent) Close() {
	for path, server := range s.servers {
		log.Println("unmounting", path)
		err := server.Unmount()
		if err != nil {
			log.Println("unmount failed", path, err)
		}
	}
}
