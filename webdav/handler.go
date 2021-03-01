package webdav

import (
	"github.com/distribyted/distribyted/fs"
	"golang.org/x/net/webdav"
)

func newHandler(fss map[string][]fs.Filesystem) *webdav.Handler {
	return &webdav.Handler{
		Prefix:     "/",
		FileSystem: newFS(fss),
		LockSystem: webdav.NewMemLS(),
	}
}
