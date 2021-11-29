package webdav

import (
	"net/http"

	"github.com/distribyted/distribyted/fs"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/webdav"
)

func newHandler(fs fs.Filesystem) *webdav.Handler {
	l := log.Logger.With().Str("component", "webDAV").Logger()
	return &webdav.Handler{
		Prefix:     "/",
		FileSystem: newFS(fs),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			if err != nil {
				l.Error().Err(err).Str("path", req.RequestURI).Msg("webDAV error")
			}
		},
	}
}
