package factory

import (
	"context"
	"print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
	"github.com/mondegor/go-webcore/mrserver/mrresponse"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"
)

func NewRestRouter(ctx context.Context, cfg config.Config, translator *mrlang.Translator) (*mrjulienrouter.RouterAdapter, error) {
	logger := mrlog.Ctx(ctx)

	corsOptions := mrrscors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		Logger:           logger.With().Str("middleware", "cors").Logger(),
	}

	errorSender, err := NewErrorResponseSender(ctx, cfg)

	if err != nil {
		return nil, err
	}

	handler, err := mrserver.NewMiddlewareHttpHandlerAdapter(errorSender)

	if err != nil {
		return nil, err
	}

	router := mrjulienrouter.New(
		logger.With().Str("router", "julienrouter").Logger(),
		handler,
		mrresponse.HandlerGetNotFoundAsJson(),
		mrresponse.HandlerGetMethodNotAllowedAsJson(),
	)
	router.RegisterMiddleware(
		mrrscors.New(corsOptions),
		mrserver.MiddlewareGeneral(translator),
	)

	return router, nil
}
