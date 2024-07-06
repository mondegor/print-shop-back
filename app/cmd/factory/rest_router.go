package factory

import (
	"context"
	"time"

	"github.com/mondegor/print-shop-back/internal/app"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrchi"
	"github.com/mondegor/go-webcore/mrserver/mrprometheus"
	"github.com/mondegor/go-webcore/mrserver/mrresp"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"
)

// NewRestRouter - создаёт объект mrchi.RouterAdapter.
func NewRestRouter(ctx context.Context, opts app.Options, translator *mrlang.Translator) (*mrchi.RouterAdapter, error) {
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

	observeRequest := mrprometheus.NewObserveRequest()
	opts.Prometheus.MustRegister(
		observeRequest.Collectors()...,
	)

	router := mrchi.New(
		logger.With().Str("router", "chi").Logger(),
		mrserver.MiddlewareHandlerAdapter(errorSender),
		mrresp.HandlerGetNotFoundAsJSON(opts.Cfg.Debugging.UnexpectedHttpStatus),
		mrresp.HandlerGetMethodNotAllowedAsJSON(opts.Cfg.Debugging.UnexpectedHttpStatus),
	)

	router.RegisterMiddleware(
		mrrscors.Middleware(corsOptions),
		mrserver.MiddlewareGeneral(
			translator,
			func(l mrlog.Logger, start time.Time, sr *mrserver.StatRequest, sw *mrserver.StatResponseWriter) {
				observeRequest.SendMetrics(l, start, sr, sw)
				mrresp.TraceRequest(l, start, sr, sw)
			},
		),
		mrserver.MiddlewareRecoverHandler(
			opts.Cfg.Debugging.Debug,
			mrresp.HandlerGetFatalErrorAsJSON(opts.Cfg.Debugging.UnexpectedHttpStatus),
		),
	)

	return router, nil
}
