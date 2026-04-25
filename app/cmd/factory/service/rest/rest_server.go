package rest

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/httpserver"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	restServerCaption = "HttpServer"
)

// InitRestServer - создаёт объект mrserver.ServerAdapter.
func InitRestServer(opts app.Options) (*httpserver.Adapter, error) {
	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", restServerCaption), "port", opts.Cfg.HttpServerPort)

	router, err := InitRestRouterWithHandlers(opts)
	if err != nil {
		return nil, err
	}

	return httpserver.New(
		router,
		httpserver.WithLogger(opts.Logger),
		httpserver.WithCaption(restServerCaption),
		httpserver.WithHostPort(opts.Cfg.HttpServerBindIP, opts.Cfg.HttpServerPort),
		httpserver.WithReadTimeout(opts.Cfg.HttpServerReadTimeout),
		httpserver.WithWriteTimeout(opts.Cfg.HttpServerWriteTimeout),
		httpserver.WithShutdownTimeout(opts.Cfg.HttpServerShutdownTimeout),
	), nil
}
