//go:build nofuse || !cgo

package fuse

import (
	"errors"

	"github.com/distribyted/distribyted/fs"
	"github.com/rs/zerolog/log"
)

// Supported indicates whether FUSE is supported in this build.
const Supported = false

// ErrFuseNotSupported is returned when attempting to mount FUSE in a build without FUSE support.
var ErrFuseNotSupported = errors.New("fuse is not supported in this build")

type Handler struct {
	fuseAllowOther bool
	path           string
}

func NewHandler(fuseAllowOther bool, path string) *Handler {
	return &Handler{
		fuseAllowOther: fuseAllowOther,
		path:           path,
	}
}

func (s *Handler) Mount(fss map[string]fs.Filesystem) error {
	log.Warn().Str("path", s.path).Msg("FUSE is disabled in this build")
	return ErrFuseNotSupported
}

func (s *Handler) Unmount() {
}
