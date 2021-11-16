package http

import (
	"net/http"

	"github.com/distribyted/distribyted/torrent"
	"github.com/gin-gonic/gin"
)

var indexHandler = func(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

var routesHandler = func(ss *torrent.Stats) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "routes.html", ss.RoutesStats())
	}
}
