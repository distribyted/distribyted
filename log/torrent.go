package log

import (
	"github.com/anacrolix/log"
	"github.com/rs/zerolog"
)

var _ log.Handler = &Torrent{}

type Torrent struct {
	L zerolog.Logger
}

func (l *Torrent) Handle(r log.Record) {
	e := l.L.Info()
	switch r.Level {
	case log.Debug:
		e = l.L.Debug()
	case log.Info:
		e = l.L.Debug().Str("error-type", "info")
	case log.Warning:
		e = l.L.Warn()
	case log.Error:
		e = l.L.Warn().Str("error-type", "error")
	case log.Critical:
		e = l.L.Warn().Str("error-type", "critical")
	}

	// TODO set log values somehow

	e.Msgf(r.Text())
}
