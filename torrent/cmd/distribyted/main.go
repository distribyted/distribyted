package main

import (
	"io/ioutil"
	"log"
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
	"github.com/panjf2000/ants/v2"
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

	pool, err := ants.NewPool(100)
	if err != nil {
		log.Fatal(err)
	}

	ss := stats.NewTorrent()
	mountService := mount.NewTorrent(c, pool, ss)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Closing torrent client...")
		c.Close()
		log.Println("Releasing execution pool...")
		pool.Release()
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
	fs := distribyted.NewBinaryFileSystem()
	file, err := fs.Open("index.html")
	if err != nil {
		log.Println("PUES SI QUE NO ESTá", err)
	} else {
		log.Println("FILE", file)
	}

	r.Use(static.Serve("/", fs))

	r.GET("/api/status", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"cacheItems":    fc.Info().NumItems,
			"cacheFilled":   fc.Info().Filled / 1024 / 1024,
			"cacheCapacity": fc.Info().Capacity / 1024 / 1024,
			"poolCap":       pool.Cap(),
			"poolFree":      pool.Free(),
			"torrentStats":  ss.Global(),
		})
	})

	r.GET("/api/status/:torrent", func(ctx *gin.Context) {
		hash := ctx.Param("torrent")
		stats, err := ss.Torrent(hash)
		if err != nil {
			ctx.AbortWithError(404, err)
		}

		ctx.JSON(200, stats)
	})

	if err := r.Run(":4444"); err != nil {
		log.Fatal(err)
	}
}
