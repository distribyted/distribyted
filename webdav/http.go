package webdav

import (
	"fmt"
	"net/http"

	"github.com/distribyted/distribyted/fs"
	"github.com/sirupsen/logrus"
)

func NewWebDAVServer(fss map[string][]fs.Filesystem, port int) error {
	logrus.WithField("host", fmt.Sprintf("0.0.0.0:%d", port)).Info("starting webDAV server")
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), newHandler(fss))
}
