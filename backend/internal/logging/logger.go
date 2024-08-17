package logging

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	zLogger *zerolog.Logger
}

func NewLogger() *Logger {
	// TODO: do not hard-code log file
	f, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't open logging file")
	}

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }

	l := zerolog.New(zerolog.MultiLevelWriter(f, os.Stdout)).With().Timestamp().Logger()

	return &Logger{
		zLogger: &l,
	}
}

func (l *Logger) NewEvent(level zerolog.Level) *zerolog.Event {
	return l.zLogger.WithLevel(level)
}
