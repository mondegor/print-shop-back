package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/initing"
)

// RegisterRestRouterCustHandlers - регистрирует в указанном роутере обработчики секции CustomersAPI.
func RegisterRestRouterCustHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	sect *section.RoutingSection,
	memberProvider mraccess.MemberProvider,
) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))
	prepareHandler := mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider)

	controllers, err := initing.CreateHttpControllers(opts.Logger, getCustomersAPIControllers(opts), prepareHandler)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getCustomersAPIControllers(_ app.Options) []initing.HttpModule {
	return nil
}
