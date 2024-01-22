package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/mondegor/go-webcore/mrserver/mrjulienrouter"
	"github.com/mondegor/go-webcore/mrserver/mrrscors"
)

func NewHttpRouter(cfg *config.Config, logger mrcore.Logger, translator *mrlang.Translator) (*mrjulienrouter.RouterAdapter, error) {
	logger.Info("Create and init http router")

	corsOptions := mrrscors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		Logger:           logger,
	}

	errorSender, err := NewErrorResponseSender(cfg, logger)

	if err != nil {
		return nil, err
	}

	handler, err := mrserver.HandlerAdapter(errorSender)

	if err != nil {
		return nil, err
	}

	router := mrjulienrouter.New(logger, handler)
	router.RegisterMiddleware(
		mrrscors.New(corsOptions),
		mrserver.MiddlewareFirst(logger, translator),
	)

	return router, nil
}
