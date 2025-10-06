package rest

import (
	"net/http"

	"github.com/mondegor/go-webcore/mraccess"
	"github.com/mondegor/go-webcore/mraccess/section"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
	"github.com/mondegor/print-shop-back/internal/factory/auth"
	authpub "github.com/mondegor/print-shop-back/internal/factory/auth/section/pub"
	"github.com/mondegor/print-shop-back/internal/initing"
)

// RegisterRestRouterAuthHandlers - регистрирует в указанном роутере обработчики секции AuthAPI.
func RegisterRestRouterAuthHandlers(
	router mrserver.HttpRouter,
	opts app.Options,
	sect *section.RoutingSection,
	memberProvider mraccess.MemberProvider,
) error {
	router.HandlerFunc(http.MethodGet, sect.BuildPath("/"), mrresp.HandlerGetStatusOkAsJSON(opts.Logger))

	controllers, err := initing.CreateHttpControllers(
		opts.Logger,
		getAuthAPIControllers(opts),
		mrinit.WithMiddlewareCheckAccess(opts.Logger, sect, memberProvider, opts.RealmKindRights, opts.PermsProvider),
	)
	if err != nil {
		return err
	}

	router.Register(controllers...)

	return nil
}

func getAuthAPIControllers(opts app.Options) []initing.HttpModule {
	return []initing.HttpModule{
		authpub.InitHttpModule(
			opts.Logger,
			opts.EventEmitter,
			opts.UseCaseErrorWrapper,
			opts.StorageErrorWrapper,
			opts.PostgresConnManager,
			opts.Locker,
			opts.RequestParsers.Parser,
			opts.ResponseSenders.Sender,
			opts.ResponseSenders.FileSender,
			opts.NotifierAPI,
			opts.Cfg.Debugging.Debug,
			opts.Cfg.AccessControl.Realms,
			opts.Cfg.AccessControl.OperationConfirm,
			auth.JWTConfig{
				Method: opts.Cfg.AccessControl.JWTMethod,
				Secret: []byte(opts.Cfg.AccessControl.JWTSecret),
			},
		),
	}
}
