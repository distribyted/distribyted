package http

import (
	"fmt"
	"net/http"

	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/httpfs/html/vfstemplate"

	"github.com/distribyted/distribyted"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/torrent"
)

func New(fc *filecache.Cache, ss *torrent.Stats, s *torrent.Service, ch *config.Handler, tss []*torrent.Server, fs http.FileSystem, logPath string, cfg *config.HTTPGlobal) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.ErrorLogger())
	r.Use(Logger())

	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(distribyted.Assets))
	})

	if cfg.HTTPFS {
		log.Info().Str("host", fmt.Sprintf("%s:%d/fs", cfg.IP, cfg.Port)).Msg("starting HTTPFS")
		h := func(c *gin.Context) {
			path := c.Param("filepath")
			c.FileFromFS(path, fs)
		}
		r.GET("/fs/*filepath", h)
		r.HEAD("/fs/*filepath", h)

	}

	t, err := vfstemplate.ParseGlob(http.FS(distribyted.Templates), nil, "/templates/*")
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

	log.Info().Str("host", fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)).Msg("starting webserver")

	if err := r.Run(fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)); err != nil {
		return fmt.Errorf("error initializing server: %w", err)
	}

	return nil
}

func Logger() gin.HandlerFunc {
	l := log.Logger.With().Str("component", "http").Logger()
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		if raw != "" {
			path = path + "?" + raw
		}
		msg := c.Errors.String()
		if msg == "" {
			msg = "Request"
		}

		s := c.Writer.Status()
		switch {
		case s >= 400 && s < 500:
			l.Warn().Str("path", path).Int("status", s).Msg(msg)
		case s >= 500:
			l.Error().Str("path", path).Int("status", s).Msg(msg)
		default:
			l.Debug().Str("path", path).Int("status", s).Msg(msg)
		}
	}
}
