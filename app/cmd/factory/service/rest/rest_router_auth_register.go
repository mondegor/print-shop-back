package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	authpub "github.com/mondegor/print-shop-back/internal/factory/auth/section/pub"
)

// RegisterRestRouterAuthHandlers - регистрирует в указанном роутере обработчики секции AuthAPI.
func RegisterRestRouterAuthHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	sect *section.RoutingSection,
	memberProvider mraccess.MemberProvider,
) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	for _, createFunc := range getAuthAPIControllers(opts) {
		list, err := createFunc()
		if err != nil {
			return err
		}

		router.Register(
			mrinit.PrepareEachController(list, prepareHandler)...,
		)
	}

	return nil
}

func getAuthAPIControllers(opts app.Options) []func() (list []mrserver.HttpController, err error) {
	return []func() (list []mrserver.HttpController, err error){
		func() ([]mrserver.HttpController, error) {
			return authpub.CreateModule(opts.AuthModule)
		},
	}
}
