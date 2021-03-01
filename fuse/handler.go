package fuse

import (
	"os"
	"path"
	"runtime"

	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fs"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	fuseAllowOther bool

	hosts map[string]*fuse.FileSystemHost
}

func NewHandler(fuseAllowOther bool) *Handler {
	return &Handler{
		fuseAllowOther: fuseAllowOther,
		hosts:          make(map[string]*fuse.FileSystemHost),
	}
}

func (s *Handler) MountAll(fss map[string][]fs.Filesystem, ef config.EventFunc) error {
	for p, fss := range fss {
		folder := p
		// On windows, the folder must don't exist
		if runtime.GOOS == "windows" {
			folder = path.Dir(folder)
		}
		if err := os.MkdirAll(folder, 0744); err != nil && !os.IsExist(err) {
			return err
		}

		host := fuse.NewFileSystemHost(NewFS(fss))

		// TODO improve error handling here
		go func() {
			var config []string

			if s.fuseAllowOther {
				config = append(config, "-o", "allow_other")
			}

			ok := host.Mount(p, config)
			if !ok {
				log.WithField("path", p).Error("error trying to mount filesystem")
			}
		}()

		s.hosts[p] = host
	}

	return nil
}

func (s *Handler) UnmountAll() {
	for path, server := range s.hosts {
		log.WithField("path", path).Info("unmounting")
		ok := server.Unmount()
		if !ok {
			//TODO try to force unmount if possible
			log.WithField("path", path).Error("unmount failed")
		}
	}

	s.hosts = make(map[string]*fuse.FileSystemHost)
}
