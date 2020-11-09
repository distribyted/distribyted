package http

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/gin-gonic/gin"
)

var apiStatusHandler = func(fc *filecache.Cache, ss *stats.Torrent) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO move to a struct
		ctx.JSON(http.StatusOK, gin.H{
			"cacheItems":    fc.Info().NumItems,
			"cacheFilled":   fc.Info().Filled / 1024 / 1024,
			"cacheCapacity": fc.Info().Capacity / 1024 / 1024,
			"torrentStats":  ss.GlobalStats(),
		})
	}
}

var apiRoutesHandler = func(ss *stats.Torrent) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		stats := ss.RoutesStats()
		ctx.JSON(http.StatusOK, stats)
	}
}

var apiGetConfigFile = func(ch *config.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rc, err := ch.GetRaw()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Data(http.StatusOK, "text/x-yaml", rc)
	}
}

var apiSetConfigFile = func(ch *config.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("no config file sent"))
			return
		}

		c, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		if len(c) == 0 {
			ctx.AbortWithError(http.StatusInternalServerError, errors.New("no config file sent"))
			return
		}

		if err := ch.Set(c); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "config file saved",
		})
	}
}

var apiStreamEvents = func(events chan string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Stream(func(w io.Writer) bool {
			if msg, ok := <-events; ok {
				ctx.SSEvent("event", msg)
				return true
			}
			return false
		})
	}
}

var apiReloadServer = func(ch *config.Handler, events chan string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		events <- "starting reload configuration process..."

		err := ch.Reload(
			func(m string) {
				events <- m
			})
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "reload process finished with no errors",
		})
	}
}
