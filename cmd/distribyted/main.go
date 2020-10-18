package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajnavarro/distribyted"
	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/fuse"
	"github.com/ajnavarro/distribyted/stats"
	tlog "github.com/anacrolix/log"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/sirupsen/logrus"
)

func main() {

	var configPath string
	if len(os.Args) < 2 {
		configPath = "./config.yaml"
	} else {
		configPath = os.Args[1]
	}

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetLevel(logrus.InfoLevel)

	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.WithError(err).Error("error reading configuration file")
		return
	}

	conf := &config.Root{}
	if err := yaml.Unmarshal(f, conf); err != nil {
		log.WithError(err).Error("error parsing configuration file")
		return
	}

	conf = config.AddDefaults(conf)

	if err := os.MkdirAll(conf.MetadataFolder, 0770); err != nil {
		log.WithError(err).Error("error creating metadata folder")
		return
	}

	fc, err := filecache.NewCache(conf.MetadataFolder)
	if err != nil {
		log.WithError(err).Error("error creating cache")
		return
	}

	fc.SetCapacity(conf.MaxCacheSize * 1024 * 1024)
	st := storage.NewResourcePieces(fc.AsResourceProvider())

	// TODO download and upload limits
	torrentCfg := torrent.NewDefaultClientConfig()
	torrentCfg.Logger = tlog.Discard
	torrentCfg.Seed = true
	torrentCfg.DisableTCP = true
	torrentCfg.DefaultStorage = st

	c, err := torrent.NewClient(torrentCfg)
	if err != nil {
		log.WithError(err).Error("error initializing torrent client")
		return
	}

	ss := stats.NewTorrent()
	mountService := fuse.NewHandler(c, ss)

	defer func() {
		tryClose(log, c, mountService)
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		tryClose(log, c, mountService)
	}()

	for _, mp := range conf.MountPoints {
		if err := mountService.Mount(mp); err != nil {
			log.WithError(err).WithField("path", mp.Path).Error("error mounting folder")
			return
		}
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(gin.Recovery())

	assets := distribyted.NewBinaryFileSystem(distribyted.HttpFS, "/assets")
	r.Use(static.Serve("/assets", assets))
	t, err := vfstemplate.ParseGlob(distribyted.HttpFS, nil, "/templates/*")
	if err != nil {
		log.WithError(err).Error("error parsing html template")
	}

	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/routes", func(c *gin.Context) {
		c.HTML(http.StatusOK, "routes.html", ss.RoutesStats())
	})

	r.GET("/api/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"cacheItems":    fc.Info().NumItems,
			"cacheFilled":   fc.Info().Filled / 1024 / 1024,
			"cacheCapacity": fc.Info().Capacity / 1024 / 1024,
			"torrentStats":  ss.GlobalStats(),
		})
	})

	r.GET("/api/routes", func(ctx *gin.Context) {
		stats := ss.RoutesStats()
		ctx.JSON(200, stats)
	})

	log.WithField("host", "0.0.0.0:4444").Info("starting webserver")

	//TODO add port from configuration
	if err := r.Run(":4444"); err != nil {
		log.WithError(err).Error("error initializing server")
		return
	}

}

func tryClose(log *logrus.Logger, c *torrent.Client, mountService *fuse.Handler) {
	log.Info("closing torrent client...")
	c.Close()
	log.Info("unmounting fuse filesystem...")
	mountService.Close()

	log.Info("exiting")
	os.Exit(1)
}
