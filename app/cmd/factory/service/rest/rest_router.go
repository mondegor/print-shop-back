package rest

import (
	"github.com/mondegor/go-components/mrauth/component/produce"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"
	"github.com/mondegor/go-webcore/mrserver/stat"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitRestRouterWithHandlers - создаёт объект mrchi.RouterAdapter и регистрирует в нём http обработчики.
func InitRestRouterWithHandlers(opts app.Options) (*mrchi.RouterAdapter, error) {
	name2section := initRoutingSections(opts)
	realm2provider := initMemberProviders(opts)

	router, err := initRestRouter(opts)
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterAuthHandlers(router, opts, name2section["auth-api"], realm2provider["*"])
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterAdmHandlers(router, opts, name2section["admin-api"], realm2provider["admin.printshop/backend"])
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterCustHandlers(router, opts, name2section["customer-api"], realm2provider["printshop/customers"])
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterProvHandlers(router, opts, name2section["provider-api"], realm2provider["printshop/providers"])
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterPubHandlers(router, opts, name2section["public-api"], realm2provider["printshop/*"])
	if err != nil {
		return nil, err
	}

	return router, nil
}

func initRestRouter(opts app.Options) (*mrchi.RouterAdapter, error) {
	corsOptions := mrrscors.Options{
		AllowedOrigins:   opts.Cfg.Cors.AllowedOrigins,
		AllowedMethods:   opts.Cfg.Cors.AllowedMethods,
		AllowedHeaders:   opts.Cfg.Cors.AllowedHeaders,
		ExposedHeaders:   opts.Cfg.Cors.ExposedHeaders,
		AllowCredentials: opts.Cfg.Cors.AllowCredentials,
		Logger:           opts.Logger.WithAttrs("middleware", "cors"),
	}

	errorSender, err := NewErrorResponseSender(opts)
	if err != nil {
		return nil, err
	}

	requestMetrics := mrprometheus.NewObserveRequest("rest_api", "go")
	opts.Prometheus.Add(requestMetrics.Collectors()...)

	requestStat := mrserver.NewRequestContainer(
		stat.NewRequestMetrics(requestMetrics),
		stat.NewRequestTracer(opts.Tracer),
		stat.NewRequestLogger(opts.Logger),
		produce.NewUserRequest(
			opts.UserStatRequestCollectorService, // TODO: заменить на API
			opts.Logger,
			opts.RequestParsers.ClientIP,
			opts.RequestParsers.User,
		),
	)

	router := mrchi.New(
		opts.Logger.WithAttrs("router", "chi"),
		mrserver.MiddlewareHandlerAdapter(errorSender),
		mrresp.HandlerGetNotFoundAsJSON(opts.Logger),
		mrresp.HandlerGetMethodNotAllowedAsJSON(opts.Logger),
	)

	router.RegisterMiddleware(
		mrserver.MiddlewareRecoverHandler(
			opts.Logger,
			opts.Cfg.Debugging.Debug,
			mrresp.HandlerGetFatalErrorAsJSON(opts.Logger),
		),
		mrserver.MiddlewareRequestID(opts.Logger, opts.TraceManager),
		mrrscors.Middleware(corsOptions),
		mrserver.MiddlewareObserver(opts.Logger, requestStat),
	)

	return router, nil
}
