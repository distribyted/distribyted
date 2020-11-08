package http

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/gin-gonic/gin"
)

var apiStatusHandler = func(fc *filecache.Cache, ss *stats.Torrent) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO move to a struct
		ctx.JSON(200, gin.H{
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
		ctx.JSON(200, stats)
	}
}

var apiGetConfigFile = func(ch *config.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rc, err := ch.GetRaw()
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		ctx.Data(200, "text/x-yaml", rc)
	}
}

var apiSetConfigFile = func(ch *config.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.AbortWithError(500, errors.New("no config file sent"))
			return
		}

		c, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		if len(c) == 0 {
			ctx.AbortWithError(500, errors.New("no config file sent"))
			return
		}

		if err := ch.Set(c); err != nil {
			ctx.AbortWithError(500, err)
			return
		}

		// TODO return something?
		ctx.JSON(200, nil)
	}
}

var apiReloadServer = func(ch *config.Handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Stream(func(w io.Writer) bool {
			err := ch.Reload(
				func(m string) {
					ctx.SSEvent("reload", m)
				})
			if err != nil {
				ctx.AbortWithError(500, err)
			}

			return false
		})
	}
}
