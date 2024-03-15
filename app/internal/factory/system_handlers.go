package factory

import (
	"context"
	"net/http"
	"print-shop-back/config"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrperms"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
)

func RegisterSystemHandlers(
	ctx context.Context,
	cfg config.Config,
	router mrserver.HttpRouter,
	section mrperms.AppSection,
) error {
	mrlog.Ctx(ctx).Info().Msgf("Init system handlers in section %s", section.Caption())

	router.HandlerFunc(http.MethodGet, section.Path("/"), mrserver.HandlerGetStatusOKAsJson())
	router.HandlerFunc(http.MethodGet, section.Path("/v1/health"), mrserver.HandlerGetHealth())

	systemInfoFunc, err := mrresponse.HandlerGetSystemInfoAsJson(
		mrresponse.SystemInfoConfig{
			Name:      cfg.AppName,
			Version:   cfg.AppVersion,
			StartedAt: cfg.AppStartedAt,
		},
	)

	if err != nil {
		return err
	}

	router.HandlerFunc(http.MethodGet, section.Path("/v1/system-info"), systemInfoFunc)

	return nil
}
