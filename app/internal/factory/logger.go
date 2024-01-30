package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrzerolog"
)

func NewLogger(cfg config.Config) (*mrzerolog.LoggerAdapter, error) {
	level, err := mrlog.ParseLevel(cfg.Log.Level)

	if err != nil {
		return nil, err
	}

	if cfg.Log.TimestampFormat != "" {
		cfg.Log.TimestampFormat, err = mrlog.ParseDateTimeFormat(cfg.Log.TimestampFormat)

		if err != nil {
			return nil, err
		}

		mrzerolog.SetDateTimeFormat(cfg.Log.TimestampFormat)
	}

	return mrzerolog.New(
		mrlog.Options{
			Level:           level,
			JsonFormat:      cfg.Log.JsonFormat,
			TimestampFormat: cfg.Log.TimestampFormat,
			ConsoleColor:    cfg.Log.ConsoleColor,
		},
	), nil
}
