package webdav

import (
	"fmt"
	"net/http"

	"github.com/distribyted/distribyted/fs"
	"github.com/rs/zerolog/log"
)

func NewWebDAVServer(fs fs.Filesystem, port int) error {
	log.Info().Str("host", fmt.Sprintf("0.0.0.0:%d", port)).Msg("starting webDAV server")
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), newHandler(fs))
}
