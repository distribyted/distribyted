package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajnavarro/distribyted/config"
	"github.com/ajnavarro/distribyted/fuse"
	"github.com/ajnavarro/distribyted/http"
	"github.com/ajnavarro/distribyted/stats"
	"github.com/ajnavarro/distribyted/torrent"
	"github.com/anacrolix/missinggo/v2/filecache"
	t "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	configFlag = "config"
	portFlag   = "http-port"
)

func main() {
	app := &cli.App{
		Name:  "distribyted",
		Usage: "Torrent client with on-demand file downloading as a filesystem.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFlag,
				Value:   "./distribyted-data/config.yaml",
				EnvVars: []string{"DISTRIBYTED_CONFIG"},
				Usage:   "YAML file containing distribyted configuration.",
			},
			&cli.IntFlag{
				Name:    portFlag,
				Value:   4444,
				EnvVars: []string{"DISTRIBYTED_HTTP_PORT"},
				Usage:   "http port for web interface",
			},
		},

		Action: func(c *cli.Context) error {
			err := load(c.String(configFlag), c.Int(portFlag))
			return err
		},

		HideHelpCommand: true,
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func newCache(folder string) (*filecache.Cache, error) {
	if err := os.MkdirAll(folder, 0744); err != nil {
		return nil, fmt.Errorf("error creating metadata folder: %w", err)
	}

	return filecache.NewCache(folder)
}

func load(configPath string, port int) error {
	ch := config.NewHandler(configPath)

	conf, err := ch.Get()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	fc, err := newCache(conf.MetadataFolder)
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	st := storage.NewResourcePieces(fc.AsResourceProvider())

	c, err := torrent.NewClient(st)
	if err != nil {
		return fmt.Errorf("error starting torrent client: %w", err)
	}

	ss := stats.NewTorrent()
	mountService := fuse.NewHandler(c, ss)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		tryClose(c, mountService)
	}()

	ch.OnReload(func(c *config.Root, ef config.EventFunc) error {
		ef("unmounting filesystems")
		mountService.UnmountAll()

		ef(fmt.Sprintf("setting cache size to %d MB", c.MaxCacheSize))
		fc.SetCapacity(c.MaxCacheSize * 1024 * 1024)

		for _, mp := range c.MountPoints {
			ef(fmt.Sprintf("mounting %v with %d torrents...", mp.Path, len(mp.Torrents)))
			if err := mountService.Mount(mp); err != nil {
				return fmt.Errorf("error mounting folder %v: %w", mp.Path, err)
			}
		}

		ef("all OK")

		return nil
	})

	if err := ch.Reload(nil); err != nil {
		return fmt.Errorf("error reloading configuration: %w", err)
	}

	defer func() {
		tryClose(c, mountService)
	}()

	return http.New(fc, ss, ch, port)
}

func tryClose(c *t.Client, mountService *fuse.Handler) {
	logrus.Info("closing torrent client...")
	c.Close()
	logrus.Info("unmounting fuse filesystem...")
	mountService.UnmountAll()

	logrus.Info("exiting")
	os.Exit(1)
}
