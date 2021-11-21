package app

import (
	"gin-api-server/internal/config"

	"github.com/rs/zerolog"
)

type Logger struct {
	Cfg *config.LoggerCfg
}

func (l *Logger) Setup() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	level, err := zerolog.ParseLevel(l.Cfg.Level)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(zerolog.Level(level))

	return nil
}
