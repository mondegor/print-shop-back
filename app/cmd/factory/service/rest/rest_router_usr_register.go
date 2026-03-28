package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	warehousing "github.com/mondegor/print-shop-back/internal/factory/warehousing/actiongroup/usr"
)

// RegisterRestRouterUsrHandlers - регистрирует в указанном роутере обработчики секции UserAPI.
func RegisterRestRouterUsrHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	actionGroup *mraccess.ActionGroup,
	userProvider mraccess.UserProvider,
) error {
	router.HandlerFunc(http.MethodGet, actionGroup.BasePath.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))

	controllers, err := initing.CreateHttpControllers(
		opts.Logger,
		getUserAPIControllers(opts),
		initing.WithCheckAccessMiddleware(opts.Logger, actionGroup, userProvider, opts.PermsProvider),
	)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getUserAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		warehousing.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.PostgresConnManager,
			opts.RequestParsers.ExtendParser,
			opts.ResponseSenders.Sender,
		),
	}
}
