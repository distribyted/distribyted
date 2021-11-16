package log

import (
	"strings"

	"github.com/rs/zerolog"
)

type Badger struct {
	L zerolog.Logger
}

func (l *Badger) Errorf(m string, f ...interface{}) {
	l.L.Error().Msgf(strings.ReplaceAll(m, "\n", ""), f...)
}

func (l *Badger) Warningf(m string, f ...interface{}) {
	l.L.Warn().Msgf(strings.ReplaceAll(m, "\n", ""), f...)
}

func (l *Badger) Infof(m string, f ...interface{}) {
	l.L.Info().Msgf(strings.ReplaceAll(m, "\n", ""), f...)
}

func (l *Badger) Debugf(m string, f ...interface{}) {
	l.L.Debug().Msgf(strings.ReplaceAll(m, "\n", ""), f...)
}
