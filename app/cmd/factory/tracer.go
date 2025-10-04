package factory

import (
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/mrwire"

	"github.com/mondegor/print-shop-back/config"
)

// InitTracer - создаёт и инициализирует трейсер.
func InitTracer(cfg config.Config, logger mrlog.Logger) mrtrace.Tracer {
	tracer := mrwire.InitTracer(
		mrwire.TracerConfig{
			Environment: cfg.App.Environment,
			Version:     cfg.App.Version,
			IsEnabled:   cfg.Trace.IsEnabled,
		},
		logger,
	)

	mrlog.Info(logger, "Create and init tracer")

	return tracer
}

// InitTraceContextManager - создаёт и инициализирует менеджер.
func InitTraceContextManager(_ config.Config, logger mrlog.Logger) (manager mrtrace.ContextManager, err error) {
	manager, err = mrwire.InitTraceContextManager(logger)
	if err != nil {
		return nil, err
	}

	mrlog.Info(logger, "Create and init trace context manager")

	return manager, nil
}
