package log

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/distribyted/distribyted/config"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

const FileName = "distribyted.log"

func Load(config *config.Log) {
	var writers []io.Writer

	// fix console colors on windows
	cso := colorable.NewColorableStdout()

	writers = append(writers, zerolog.ConsoleWriter{Out: cso, TimeFormat: time.DateTime})
	writers = append(writers, newRollingFile(config))
	mw := io.MultiWriter(writers...)

	log.Logger = log.Output(mw)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	l := zerolog.InfoLevel
	if config.Debug {
		l = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(l)
}

func newRollingFile(config *config.Log) io.Writer {
	if err := os.MkdirAll(config.Path, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Path).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   filepath.Join(config.Path, FileName),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
