package factory

import (
	"net/http"
	"print-shop-back/config"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

func RegisterSystemHandlers(
	cfg *config.Config,
	logger mrcore.Logger,
	router mrserver.HttpRouter,
	section mrcore.AppSection,
) error {
	logger.Info("Init system handlers in section %s", section.Caption())

	router.HandlerFunc(http.MethodGet, section.Path("/"), mrserver.HandlerGetStatusOKAsJson())
	router.HandlerFunc(http.MethodGet, section.Path("/health"), mrserver.HandlerGetHealth())

	serviceInfoFunc, err := mrserver.HandlerGetServiceInfoAsJson(
		mrserver.ConfigServiceInfo{
			Name:      cfg.AppName,
			Version:   cfg.AppVersion,
			StartedAt: cfg.AppStartedAt,
		},
	)

	if err != nil {
		return err
	}

	router.HandlerFunc(http.MethodGet, section.Path("/service-info"), serviceInfoFunc)

	return nil
}
