package factory

import (
	"strings"

	"github.com/mondegor/print-shop-back/config"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog/zerolog"
	"github.com/mondegor/go-webcore/mrlog/zerolog/factory"
)

// NewLogger - создаёт объект zerolog.LoggerAdapter.
func NewLogger(cfg config.Config) (*zerolog.LoggerAdapter, error) {
	return factory.NewZeroLogAdapter(
		factory.Options{
			Level:            strings.ToUpper(cfg.Log.Level),
			JsonFormat:       cfg.Log.JsonFormat,
			TimestampFormat:  cfg.Log.TimestampFormat,
			ConsoleColor:     cfg.Log.ConsoleColor,
			PrepareErrorFunc: mrcore.PrepareError,
		},
	)
}
