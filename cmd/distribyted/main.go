package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajnavarro/distribyted"
	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/mount"
	"github.com/ajnavarro/distribyted/stats"
	tlog "github.com/anacrolix/log"
	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
	"github.com/shurcooL/httpfs/html/vfstemplate"
)

func main() {
	var configPath string
	if len(os.Args) < 2 {
		configPath = "./config.yaml"
	} else {
		configPath = os.Args[1]
	}

	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	conf := &config.Root{}
	if err := yaml.Unmarshal(f, conf); err != nil {
		log.Fatal(err)
	}

	conf = config.AddDefaults(conf)

	if err := os.MkdirAll(conf.MetadataFolder, 0770); err != nil {
		log.Fatal(err)
	}

	fc, err := filecache.NewCache(conf.MetadataFolder)
	if err != nil {
		log.Fatal(err)
	}

	fc.SetCapacity(conf.MaxCacheSize * 1024 * 1024)
	st := storage.NewResourcePieces(fc.AsResourceProvider())

	// TODO download and upload limits
	torrentCfg := torrent.NewDefaultClientConfig()
	torrentCfg.Logger = tlog.Default.WithDefaultLevel(tlog.Info).FilterLevel(tlog.Info)
	torrentCfg.Seed = true
	torrentCfg.DisableTCP = true
	torrentCfg.DefaultStorage = st

	c, err := torrent.NewClient(torrentCfg)
	if err != nil {
		log.Fatal(err)
	}

	ss := stats.NewTorrent()
	mountService := mount.NewHandler(c, ss)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Closing torrent client...")
		c.Close()
		log.Println("Unmounting fuse filesystem...")
		mountService.Close()

		log.Println("Exiting...")
		os.Exit(1)
	}()

	for _, mp := range conf.MountPoints {
		if err := mountService.Mount(mp); err != nil {
			log.Fatal(err)
		}
	}

	r := gin.Default()
	assets := distribyted.NewBinaryFileSystem(distribyted.Assets)
	r.Use(static.Serve("/assets", assets))

	t, err := vfstemplate.ParseGlob(distribyted.Templates, nil, "*")
	if err != nil {
		log.Fatal(err)
	}

	r.SetHTMLTemplate(t)

	//	r.LoadHTMLGlob("templates/*")

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

	if err := r.Run(":4444"); err != nil {
		log.Fatal(err)
	}
}
