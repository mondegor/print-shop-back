package factory

import (
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/wire"
	"github.com/mondegor/go-sysmess/wire/slog"

	"github.com/mondegor/print-shop-back/config"
)

// InitLoggerAndTracer - создаёт и инициализирует логгер и трейсер на основе логгера.
func InitLoggerAndTracer(cfg config.Config) (mrlog.Logger, mrtrace.Tracer, error) {
	logger, err := slog.InitLogger(
		wire.LoggerConfig{
			Environment: cfg.App.Environment,
			Version:     cfg.App.Version,
			Level:       cfg.Log.Level,
			JsonFormat:  cfg.Log.JsonFormat,
			TimeFormat:  cfg.Log.TimeFormat,
			ColorMode:   cfg.Log.ColorMode,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	mrlog.Info(logger, "Create and init logger and tracer")

	return logger, slog.InitTracer(logger), nil
}

// InitTraceContextManager - создаёт и инициализирует менеджер.
func InitTraceContextManager(_ config.Config, logger mrlog.Logger) (manager mrtrace.ContextManager, err error) {
	manager, err = wire.InitTraceContextManager(wire.DefaultProcessIDs(), logger)
	if err != nil {
		return nil, err
	}

	mrlog.Info(logger, "Create and init trace context manager")

	return manager, nil
}
