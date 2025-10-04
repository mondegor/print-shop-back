package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrsentry"

	"github.com/mondegor/print-shop-back/config"
)

// InitSentry - создаёт объект mrsentry.Adapter.
func InitSentry(logger mrlog.Logger, cfg config.Config) (*mrsentry.Adapter, error) {
	mrlog.Info(logger, "Create and init sentry")

	client, err := mrsentry.New(
		mrsentry.Options{
			DSN:              cfg.Sentry.DSN,
			Environment:      cfg.App.Environment,
			AppVersion:       cfg.App.Version,
			TracesSampleRate: cfg.Sentry.TracesSampleRate,
			FlushTimeout:     cfg.Sentry.FlushTimeout,
			StackTraceBounds: cfg.Debugging.ErrorCaller.UpperBounds,
			IsDebug:          cfg.Debugging.Debug,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("sentry.Init: %w", err)
	}

	return client, nil
}
