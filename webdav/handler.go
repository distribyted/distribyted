package webdav

import (
	"github.com/distribyted/distribyted/fs"
	"golang.org/x/net/webdav"
)

func newHandler(fs fs.Filesystem) *webdav.Handler {
	return &webdav.Handler{
		Prefix:     "/",
		FileSystem: newFS(fs),
		LockSystem: webdav.NewMemLS(),
	}
}
