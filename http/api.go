package http

import (
	"net/http"
	"sort"

	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/distribyted/distribyted/torrent"
	"github.com/gin-gonic/gin"
)

var apiStatusHandler = func(fc *filecache.Cache, ss *torrent.Stats) gin.HandlerFunc {
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

var apiRoutesHandler = func(ss *torrent.Stats) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := ss.RoutesStats()
		sort.Sort(torrent.ByName(s))
		ctx.JSON(http.StatusOK, s)
	}
}

var apiAddTorrentHandler = func(s *torrent.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		route := ctx.Param("route")

		var json RouteAdd
		if err := ctx.ShouldBindJSON(&json); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.AddMagnet(route, json.Magnet); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

var apiDelTorrentHandler = func(s *torrent.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		route := ctx.Param("route")
		hash := ctx.Param("torrent_hash")

		if err := s.RemoveFromHash(route, hash); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}
