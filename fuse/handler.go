package fuse

import (
	"os"
	"path"
	"runtime"

	"github.com/billziss-gh/cgofuse/fuse"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fs"
	"github.com/rs/zerolog/log"
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

func (s *Handler) MountAll(fss map[string]fs.Filesystem, ef config.EventFunc) error {
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
				log.Error().Str("path", p).Msg("error trying to mount filesystem")
			}
		}()

		s.hosts[p] = host
	}

	return nil
}

func (s *Handler) UnmountAll() {
	for path, server := range s.hosts {
		log.Info().Str("path", path).Msg("unmounting")
		ok := server.Unmount()
		if !ok {
			//TODO try to force unmount if possible
			log.Error().Str("path", path).Msg("unmount failed")
		}
	}

	s.hosts = make(map[string]*fuse.FileSystemHost)
}
