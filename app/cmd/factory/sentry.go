package factory

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrclient/sentry"
	"github.com/pkg/errors"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
)

var errSentryDisabled = errors.New("sentry disabled")

// InitSentry - создаёт объект sentry.Adapter.
func InitSentry(logger log.Logger, cfg config.Config) (*sentry.Adapter, error) {
	if cfg.SentryDSN == "" {
		return nil, errSentryDisabled
	}

	log.Info(logger, "Create and init sentry")

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
