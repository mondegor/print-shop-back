package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrclient/sentry"
	"github.com/pkg/errors"

	"github.com/mondegor/print-shop-back/config"
)

var errSentryDisabled = errors.New("sentry disabled")

// InitSentry - создаёт объект sentry.Adapter.
func InitSentry(logger mrlog.Logger, cfg config.Config) (*sentry.Adapter, error) {
	if cfg.Sentry.DSN == "" {
		return nil, errSentryDisabled
	}

	mrlog.Info(logger, "Create and init sentry")

	client, err := sentry.New(
		cfg.Sentry.DSN,
		sentry.WithEnvironment(cfg.App.Environment),
		sentry.WithRelease(cfg.App.Version),
		sentry.WithDebugMode(cfg.Debugging.Debug),
		sentry.WithTracesSampleRate(cfg.Sentry.TracesSampleRate),
		sentry.WithFlushTimeout(cfg.Sentry.FlushTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("sentry.Init: %w", err)
	}

	return client, nil
}
