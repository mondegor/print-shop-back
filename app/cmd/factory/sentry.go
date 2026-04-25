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
	if cfg.SentryDSN == "" {
		return nil, errSentryDisabled
	}

	mrlog.Info(logger, "Create and init sentry")

	client, err := sentry.New(
		cfg.SentryDSN,
		sentry.WithEnvironment(cfg.Environment),
		sentry.WithRelease(cfg.AppVersion),
		sentry.WithDebugMode(cfg.DebugIsEnabled),
		sentry.WithTracesSampleRate(cfg.SentryTracesSampleRate),
		sentry.WithFlushTimeout(cfg.SentryFlushTimeout),
	)
	if err != nil {
		return nil, fmt.Errorf("sentry.Init: %w", err)
	}

	return client, nil
}
