package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/print-shop-back/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
)

// RegisterSystemHandlers - comment func.
func RegisterSystemHandlers(ctx context.Context, cfg config.Config, router mrserver.HttpRouter, section *mrperms.AppSection) error {
	mrlog.Ctx(ctx).Info().Msgf("Init system handlers in section %s", section.Caption())

	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())
	router.HandlerFunc(http.MethodGet, section.BuildPath("/v1/health"), mrresp.HandlerGetHealth())

	systemInfoFunc, err := mrresp.HandlerGetSystemInfoAsJSON(
		mrresp.SystemInfoConfig{
			Name:        cfg.App.Name,
			Version:     cfg.App.Version,
			Environment: cfg.App.Environment,
			IsDebug:     cfg.Debugging.Debug,
			StartedAt:   cfg.App.StartedAt,
		},
	)
	if err != nil {
		return err
	}

	router.HandlerFunc(http.MethodGet, section.BuildPath("/v1/system-info"), systemInfoFunc)

	return nil
}
