package rest

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/httpserver"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	restServerCaption = "RestServer"
)

// InitRestServer - создаёт объект mrserver.ServerAdapter.
func InitRestServer(opts app.Options) (*httpserver.Adapter, error) {
	srvOpts := opts.Cfg.Servers.RestServer

	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", restServerCaption), "port", srvOpts.Listen.Port)

	router, err := InitRestRouterWithHandlers(opts)
	if err != nil {
		return nil, err
	}

	return httpserver.New(
		router,
		httpserver.WithLogger(opts.Logger),
		httpserver.WithCaption(restServerCaption),
		httpserver.WithHostAndPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		httpserver.WithReadTimeout(srvOpts.ReadTimeout),
		httpserver.WithWriteTimeout(srvOpts.WriteTimeout),
		httpserver.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	), nil
}
