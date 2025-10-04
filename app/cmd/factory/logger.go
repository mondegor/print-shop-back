package factory

import (
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrwire"

	"github.com/mondegor/print-shop-back/config"
)

// InitLogger - создаёт и инициализирует логгер.
func InitLogger(cfg config.Config) (logger mrlog.Logger, err error) {
	logger, err = mrwire.InitLogger(
		mrwire.LoggerConfig{
			Environment: cfg.App.Environment,
			Version:     cfg.App.Version,
			Level:       cfg.Log.Level,
			JsonFormat:  cfg.Log.JsonFormat,
			TimeFormat:  cfg.Log.TimeFormat,
			ColorMode:   cfg.Log.ColorMode,
		},
	)
	if err != nil {
		return nil, err
	}

	mrlog.Info(logger, "Create and init logger")

	return logger, nil
}
