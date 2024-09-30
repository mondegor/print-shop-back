package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/go-webcore/mrfactory"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	provideraccountsprov "github.com/mondegor/print-shop-back/internal/factory/provideraccounts/section/prov"
)

// RegisterRestRouterProvHandlers - регистрирует в указанном роутере обработчики секции ProvidersAPI.
func RegisterRestRouterProvHandlers(ctx context.Context, router mrserver.HttpRouter, opts app.Options) error {
	section := NewAppSectionProvidersAPI(ctx, opts)
	prepareHandler := mrfactory.WithMiddlewareCheckAccess(ctx, section, opts.AccessControl)
	router.HandlerFunc(http.MethodGet, section.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON())

	for _, createFunc := range getProvidersAPIControllers(ctx, opts) {
		list, err := createFunc()
		if err != nil {
			return err
		}

		router.Register(
			mrfactory.PrepareEachController(list, prepareHandler)...,
		)
	}

	return nil
}

func getProvidersAPIControllers(ctx context.Context, opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return provideraccountsprov.CreateModule(ctx, opts.ProviderAccountsModule)
		},
	}
}
