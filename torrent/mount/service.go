package mount

import (
	"log"
	"os"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/node"
	"github.com/anacrolix/torrent"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/panjf2000/ants/v2"
)

type Service struct {
	c    *torrent.Client
	opts *fs.Options

	pool    *ants.Pool
	servers map[string]*fuse.Server
}

func NewService(c *torrent.Client, pool *ants.Pool) *Service {
	return &Service{
		c:       c,
		opts:    &fs.Options{},
		pool:    pool,
		servers: make(map[string]*fuse.Server),
	}
}

func (s *Service) Mount(mpc *config.MountPoint) error {
	var torrents []*torrent.Torrent
	for _, magnet := range mpc.Magnets {
		t, err := s.c.AddMagnet(magnet.URI)
		if err != nil {
			return err
		}
		log.Println("getting torrent info", t.Name())
		<-t.GotInfo()
		log.Println("torrent info obtained", t.Name())
		torrents = append(torrents, t)
	}

	// TODO change permissions
	if err := os.MkdirAll(mpc.Path, 0770); err != nil {
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

func (s *Service) Close() {
	for path, server := range s.servers {
		log.Println("unmounting", path)
		err := server.Unmount()
		if err != nil {
			log.Println("unmount failed", path, err)
		}
	}
}
