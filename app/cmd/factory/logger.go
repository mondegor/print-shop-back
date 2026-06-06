package factory

import (
	"github.com/mondegor/go-sysmess/wire/mrlog"
	"github.com/mondegor/go-sysmess/wire/mrlog/slog"
	"github.com/mondegor/go-sysmess/wire/mrtrace"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
)

// InitLoggerAndTracer - создаёт и инициализирует логгер и трейсер на основе логгера.
func InitLoggerAndTracer(cfg config.Config) (log.Logger, trace.Tracer, error) {
	logger, err := slog.InitLogger(
		mrlog.LoggerConfig{
			Environment:       cfg.Environment,
			Version:           cfg.AppVersion,
			Level:             cfg.LogLevel,
			JsonFormat:        cfg.LogJsonFormat,
			TimeFormat:        cfg.LogTimeFormat,
			ColorMode:         cfg.LogColorMode,
			ContextProcessIDs: mrlog.DefaultProcessIDs(),
		},
	)
	if err != nil {
		return nil, nil, err
	}

	log.Info(logger, "Create and init logger and tracer")

	return logger, slog.InitTracer(logger), nil
}

// InitTraceContextManager - создаёт и инициализирует менеджер.
func InitTraceContextManager(_ config.Config, logger log.Logger) (manager trace.ContextManager, err error) {
	manager, err = mrtrace.InitTraceContextManager(mrlog.DefaultProcessIDs(), logger)
	if err != nil {
		return nil, err
	}

	log.Info(logger, "Create and init trace context manager")

	return manager, nil
}
