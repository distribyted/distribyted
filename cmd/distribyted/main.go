package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/anacrolix/missinggo/v2/filecache"
	"github.com/anacrolix/torrent/storage"
	"github.com/distribyted/distribyted/config"
	"github.com/distribyted/distribyted/fs"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/distribyted/distribyted/fuse"
	"github.com/distribyted/distribyted/http"
	dlog "github.com/distribyted/distribyted/log"
	"github.com/distribyted/distribyted/torrent"
	"github.com/distribyted/distribyted/torrent/loader"
	"github.com/distribyted/distribyted/webdav"
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

			// stop program execution on errors to avoid flashing consoles
			if err != nil && runtime.GOOS == "windows" {
				log.Error().Err(err).Msg("problem starting application")
				fmt.Print("Press 'Enter' to continue...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
			}

			return err
		},

		HideHelpCommand: true,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("problem starting application")
	}
}

func load(configPath string, port, webDAVPort int, fuseAllowOther bool) error {
	ch := config.NewHandler(configPath)

	conf, err := ch.Get()
	if err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	dlog.Load(conf.Log)

	if err := os.MkdirAll(conf.Torrent.MetadataFolder, 0744); err != nil {
		return fmt.Errorf("error creating metadata folder: %w", err)
	}

	cf := filepath.Join(conf.Torrent.MetadataFolder, "cache")
	fc, err := filecache.NewCache(cf)
	if err != nil {
		return fmt.Errorf("error creating cache: %w", err)
	}

	st := storage.NewResourcePieces(fc.AsResourceProvider())

	// cache is not working with windows
	if runtime.GOOS == "windows" {
		st = storage.NewFile(cf)
	}

	fis, err := torrent.NewFileItemStore(filepath.Join(conf.Torrent.MetadataFolder, "items"), 2*time.Hour)
	if err != nil {
		return fmt.Errorf("error starting item store: %w", err)
	}

	c, err := torrent.NewClient(st, fis, conf.Torrent)
	if err != nil {
		return fmt.Errorf("error starting torrent client: %w", err)
	}

	pcp := filepath.Join(conf.Torrent.MetadataFolder, "piece-completion")
	if err := os.MkdirAll(pcp, 0744); err != nil {
		return fmt.Errorf("error creating piece completion folder: %w", err)
	}

	pc, err := storage.NewSqlitePieceCompletion(pcp)
	if err != nil {
		return fmt.Errorf("error creating servers piece completion: %w", err)
	}

	var servers []*torrent.Server
	for _, s := range conf.Servers {
		server := torrent.NewServer(c, pc, s)
		servers = append(servers, server)
		if err := server.Start(); err != nil {
			return fmt.Errorf("error starting server: %w", err)
		}
	}

	cl := loader.NewConfig(conf.Routes)
	ss := torrent.NewStats()

	dbl, err := loader.NewDB(filepath.Join(conf.Torrent.MetadataFolder, "magnetdb"))
	if err != nil {
		return fmt.Errorf("error starting magnet database: %w", err)
	}

	ts := torrent.NewService(cl, dbl, ss, c, conf.Torrent.AddTimeout)

	mh := fuse.NewHandler(fuseAllowOther || conf.Fuse.AllowOther, conf.Fuse.Path)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		<-sigChan
		log.Info().Msg("closing servers...")
		for _, s := range servers {
			if err := s.Close(); err != nil {
				log.Warn().Err(err).Msg("problem closing server")
			}
		}
		log.Info().Msg("closing items database...")
		fis.Close()
		log.Info().Msg("closing magnet database...")
		dbl.Close()
		log.Info().Msg("closing torrent client...")
		c.Close()
		log.Info().Msg("unmounting fuse filesystem...")
		mh.Unmount()

		log.Info().Msg("exiting")
		os.Exit(1)
	}()

	log.Info().Msg(fmt.Sprintf("setting cache size to %d MB", conf.Torrent.GlobalCacheSize))
	fc.SetCapacity(conf.Torrent.GlobalCacheSize * 1024 * 1024)

	fss, err := ts.Load()
	if err != nil {
		return fmt.Errorf("error when loading torrents: %w", err)
	}

	go func() {
		if err := mh.Mount(fss); err != nil {
			log.Info().Err(err).Msg("error mounting filesystems")
		}
	}()

	go func() {
		if conf.WebDAV != nil {
			port = webDAVPort
			if port == 0 {
				port = conf.WebDAV.Port
			}

			cfs, err := fs.NewContainerFs(fss)
			if err != nil {
				log.Error().Err(err).Msg("error adding files to webDAV")
				return
			}

			if err := webdav.NewWebDAVServer(cfs, port, conf.WebDAV.User, conf.WebDAV.Pass); err != nil {
				log.Error().Err(err).Msg("error starting webDAV")
			}
		}

		log.Warn().Msg("webDAV configuration not found!")
	}()

	logFilename := filepath.Join(conf.Log.Path, dlog.FileName)

	err = http.New(fc, ss, ts, ch, servers, port, logFilename)
	log.Error().Err(err).Msg("error initializing HTTP server")
	return err
}
