package factory

import (
	"context"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/zerolog"
	"github.com/mondegor/go-webcore/mrlog/zerolog/factory"

	"github.com/mondegor/print-shop-back/config"
)

// InitLogger - инициализирует логгер указанный в ctx,
// если он не был создан, то предварительно создаёт его и регистрирует в контексте.
func InitLogger(ctx context.Context, cfg config.Config) (context.Context, error) {
	if mrlog.HasCtx(ctx) {
		mrlog.Ctx(ctx).Info().Msg("Init logger")

		return ctx, nil
	}

	logger, err := NewZeroLogger(cfg)
	if err != nil {
		return nil, err
	}

	logger.Info().Msg("Create and init Zero logger")

	return mrlog.WithContext(ctx, logger), nil
}

// NewZeroLogger - создаёт объект zerolog.LoggerAdapter.
func NewZeroLogger(cfg config.Config) (*zerolog.LoggerAdapter, error) {
	return factory.NewZeroLogAdapter(
		factory.Options{
			Stdout:           cfg.Os.Stdout,
			Level:            strings.ToUpper(cfg.Log.Level),
			JsonFormat:       cfg.Log.JsonFormat,
			TimestampFormat:  cfg.Log.TimestampFormat,
			ConsoleColor:     cfg.Log.ConsoleColor,
			PrepareErrorFunc: mrcore.PrepareError,
		},
	)
}
