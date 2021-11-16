package log

import (
	"github.com/anacrolix/log"
	"github.com/rs/zerolog"
)

var _ log.LoggerImpl = &Torrent{}

type Torrent struct {
	L zerolog.Logger
}

func (l *Torrent) Log(m log.Msg) {
	level, ok := m.GetLevel()

	e := l.L.Info()

	if !ok {
		level = log.Debug
	}

	switch level {
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
	case log.Fatal:
		e = l.L.Warn().Str("error-type", "fatal")
	}

	e.Msgf(m.String())
}
