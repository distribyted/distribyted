package http

import (
	"net/http"

	"github.com/distribyted/distribyted/stats"
	"github.com/gin-gonic/gin"
)

var indexHandler = func(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

var routesHandler = func(ss *stats.Torrent) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "routes.html", ss.RoutesStats())
	}
}

var configHandler = func(c *gin.Context) {
	c.HTML(http.StatusOK, "config.html", nil)
}
