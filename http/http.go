package http

import (
	"fmt"

	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/distribyted/distribyted"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/torrent"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

func New(fc *filecache.Cache, ss *torrent.Stats, s *torrent.Service, ch *config.Handler, tss []*torrent.Server, port int, logPath string) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.ErrorLogger())

	assets := distribyted.NewBinaryFileSystem(distribyted.HttpFS, "/assets")
	r.Use(static.Serve("/assets", assets))
	t, err := vfstemplate.ParseGlob(distribyted.HttpFS, nil, "/templates/*")
	if err != nil {
		return fmt.Errorf("error parsing html: %w", err)
	}

	r.SetHTMLTemplate(t)

	r.GET("/", indexHandler)
	r.GET("/routes", routesHandler(ss))
	r.GET("/logs", logsHandler)
	r.GET("/servers", serversFoldersHandler())

	api := r.Group("/api")
	{
		api.GET("/log", apiLogHandler(logPath))
		api.GET("/status", apiStatusHandler(fc, ss))
		api.GET("/servers", apiServersHandler(tss))

		api.GET("/routes", apiRoutesHandler(ss))
		api.POST("/routes/:route/torrent", apiAddTorrentHandler(s))
		api.DELETE("/routes/:route/torrent/:torrent_hash", apiDelTorrentHandler(s))

	}

	log.Info().Str("host", fmt.Sprintf("0.0.0.0:%d", port)).Msg("starting webserver")

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		return fmt.Errorf("error initializing server: %w", err)
	}

	return nil
}
