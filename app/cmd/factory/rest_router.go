package factory

import (
	"context"
	"strconv"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewRestRouter - создаёт объект mrchi.RouterAdapter.
func NewRestRouter(ctx context.Context, opts app.Options) (*mrchi.RouterAdapter, error) {
	logger := mrlog.Ctx(ctx)

	corsOptions := mrrscors.Options{
		AllowedOrigins:   opts.Cfg.Cors.AllowedOrigins,
		AllowedMethods:   opts.Cfg.Cors.AllowedMethods,
		AllowedHeaders:   opts.Cfg.Cors.AllowedHeaders,
		ExposedHeaders:   opts.Cfg.Cors.ExposedHeaders,
		AllowCredentials: opts.Cfg.Cors.AllowCredentials,
		Logger:           logger.With().Str("middleware", "cors").Logger(),
	}

	errorSender, err := NewErrorResponseSender(ctx, opts)
	if err != nil {
		return nil, err
	}

	observeRequest := mrprometheus.NewObserveRequest("rest_api", "go")
	opts.Prometheus.MustRegister(
		observeRequest.Collectors()...,
	)

	router := mrchi.New(
		logger.With().Str("router", "chi").Logger(),
		mrserver.MiddlewareHandlerAdapter(errorSender),
		mrresp.HandlerGetNotFoundAsJSON(),
		mrresp.HandlerGetMethodNotAllowedAsJSON(),
	)

	router.RegisterMiddleware(
		mrrscors.Middleware(corsOptions),
		mrserver.MiddlewareGeneral(
			opts.Translator,
			func(l mrlog.Logger, start time.Time, sr *mrserver.StatRequestReader, sw *mrserver.StatResponseWriter) {
				method := sr.Request().Method
				location := sr.Request().URL.Path

				observeRequest.SetStatusWithTime(method, location, strconv.Itoa(sw.StatusCode()), start)
				observeRequest.IncrementRequestSize(method, location, sr.Size())
				observeRequest.IncrementResponseSize(method, location, sw.Size())

				mrserver.TraceRequest(l, start, sr, sw)
			},
		),
		mrserver.MiddlewareRecoverHandler(
			opts.Cfg.Debugging.Debug,
			mrresp.HandlerGetFatalErrorAsJSON(),
		),
	)

	return router, nil
}

// NewRestRouterWithRegisterHandlers - создаёт объект mrchi.RouterAdapter и регистрирует в нём http обработчики.
func NewRestRouterWithRegisterHandlers(ctx context.Context, opts app.Options) (*mrchi.RouterAdapter, error) {
	router, err := NewRestRouter(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterAdmHandlers(ctx, router, opts)
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterProvHandlers(ctx, router, opts)
	if err != nil {
		return nil, err
	}

	err = RegisterRestRouterPubHandlers(ctx, router, opts)
	if err != nil {
		return nil, err
	}

	return router, nil
}
