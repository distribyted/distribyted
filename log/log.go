package log

import (
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Load() {
	// fix console colors on windows
	cso := colorable.NewColorableStdout()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: cso})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
