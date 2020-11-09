package http

import (
	"fmt"

	"github.com/ajnavarro/distribyted"
	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/sirupsen/logrus"
)

func New(fc *filecache.Cache, ss *stats.Torrent, ch *config.Handler, port int) error {
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
	r.GET("/config", configHandler)

	eventChan := make(chan string)

	api := r.Group("/api")
	{
		api.GET("/status", apiStatusHandler(fc, ss))
		api.GET("/routes", apiRoutesHandler(ss))
		api.GET("/config", apiGetConfigFile(ch))
		api.POST("/config", apiSetConfigFile(ch))
		api.POST("/reload", apiReloadServer(ch, eventChan))
		api.GET("/events", apiStreamEvents(eventChan))
	}

	logrus.WithField("host", fmt.Sprintf("0.0.0.0:%d", port)).Info("starting webserver")

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		return fmt.Errorf("error initializing server: %w", err)
	}

	return nil
}
