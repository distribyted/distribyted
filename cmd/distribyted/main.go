package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/anacrolix/missinggo/v2/filecache"
	t "github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/storage"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fuse"
	"github.com/distribyted/distribyted/http"
	"github.com/distribyted/distribyted/stats"
	"github.com/distribyted/distribyted/torrent"
	"github.com/distribyted/distribyted/webdav"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	configFlag     = "config"
	fuseAllowOther = "fuse-allow-other"
	portFlag       = "http-port"
	webDAVPortFlag = "webdav-port"
)

func main() {
	app := &cli.App{
		Name:  "distribyted",
		Usage: "Torrent client with on-demand file downloading as a filesystem.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    configFlag,
				Value:   "./distribyted-data/config/config.yaml",
				EnvVars: []string{"DISTRIBYTED_CONFIG"},
				Usage:   "YAML file containing distribyted configuration.",
			},
			&cli.IntFlag{
				Name:    portFlag,
				Value:   4444,
				EnvVars: []string{"DISTRIBYTED_HTTP_PORT"},
				Usage:   "HTTP port for web interface.",
			},
			&cli.IntFlag{
				Name:    webDAVPortFlag,
				Value:   36911,
				EnvVars: []string{"DISTRIBYTED_WEBDAV_PORT"},
				Usage:   "Port used for WebDAV interface.",
			},
			&cli.BoolFlag{
				Name:    fuseAllowOther,
				Value:   false,
				EnvVars: []string{"DISTRIBYTED_FUSE_ALLOW_OTHER"},
				Usage:   "Allow other users to access all fuse mountpoints. You need to add user_allow_other flag to /etc/fuse.conf file.",
			},
		},

		Action: func(c *cli.Context) error {
			err := load(c.String(configFlag), c.Int(portFlag), c.Int(webDAVPortFlag), c.Bool(fuseAllowOther))
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

func load(configPath string, port, webDAVPort int, fuseAllowOther bool) error {
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

	th := torrent.NewHandler(c, ss)

	mh := fuse.NewHandler(fuseAllowOther || conf.AllowOther)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		tryClose(c, mh)
	}()

	ch.OnReload(func(c *config.Root, ef config.EventFunc) error {
		ef("unmounting filesystems")
		mh.UnmountAll()
		th.RemoveAll()

		ef(fmt.Sprintf("setting cache size to %d MB", c.MaxCacheSize))
		fc.SetCapacity(c.MaxCacheSize * 1024 * 1024)

		for _, mp := range c.MountPoints {
			ef(fmt.Sprintf("loading %v with %d torrents...", mp.Path, len(mp.Torrents)))
			if err := th.Load(mp.Path, mp.Torrents); err != nil {
				return fmt.Errorf("error loading folder %v: %w", mp.Path, err)
			}
			ef(fmt.Sprintf("%v loaded", mp.Path))
		}

		return mh.MountAll(th.Fileststems(), ef)

	})

	if err := ch.Reload(nil); err != nil {
		return fmt.Errorf("error reloading configuration: %w", err)
	}

	defer func() {
		tryClose(c, mh)
	}()

	go func() {
		if conf.WebDAV != nil {
			wdth := torrent.NewHandler(c, ss)
			if err := wdth.Load("::/webDAV", conf.WebDAV.Torrents); err != nil {
				logrus.WithError(err).Error("error loading torrents for webDAV")
			}

			if err := webdav.NewWebDAVServer(wdth.Fileststems(), webDAVPort); err != nil {
				logrus.WithError(err).Error("error starting webDAV")
			}
		}
	}()

	err = http.New(fc, ss, ch, port)

	logrus.WithError(err).Error("error initializing HTTP server")

	return err
}

func tryClose(c *t.Client, mountService *fuse.Handler) {
	logrus.Info("closing torrent client...")
	c.Close()
	logrus.Info("unmounting fuse filesystem...")
	mountService.UnmountAll()

	logrus.Info("exiting")
	os.Exit(1)
}
