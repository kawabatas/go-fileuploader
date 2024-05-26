package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type stdLogger struct{}

var _ Logger = (*stdLogger)(nil)

func newStdLogger() *stdLogger {
	// TODO: slog に変える
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	return &stdLogger{}
}

func (l *stdLogger) Debugf(format string, v ...any) {
	if v != nil {
		log.Debug().Msgf(format, v...)
		return
	}
	log.Debug().Msg(format)
}

func (l *stdLogger) Infof(format string, v ...any) {
	if v != nil {
		log.Info().Msgf(format, v...)
		return
	}
	log.Info().Msg(format)
}

func (l *stdLogger) Errorf(format string, v ...any) {
	if v != nil {
		log.Error().Msgf(format, v...)
		return
	}
	log.Error().Msg(format)
}
