package factory

import (
	"print-shop-back/config"

	"github.com/mondegor/go-sysmess/mrlang"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrserver"
)

func NewHttpRouter(cfg *config.Config, logger mrcore.Logger, translator *mrlang.Translator) (mrcore.HttpRouter, error) {
	requestValidator, err := NewValidator(cfg, logger)

	if err != nil {
		return nil, err
	}

	logger.Info("Create and init http router")

	corsOptions := mrserver.CorsOptions{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		Logger:           logger,
	}

	router := mrserver.NewRouter(logger, mrserver.HandlerAdapter(nil))
	router.RegisterMiddleware(
		mrserver.NewCors(corsOptions),
		mrserver.MiddlewareFirst(logger, translator, requestValidator),
	)

	return router, nil
}
