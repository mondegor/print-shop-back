package rest

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrhttp"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	restServerCaption = "RestServer"
)

// InitRestServer - создаёт объект mrserver.ServerAdapter.
func InitRestServer(opts app.Options) (*mrhttp.Adapter, error) {
	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", restServerCaption))

	router, err := InitRestRouterWithHandlers(opts)
	if err != nil {
		return nil, err
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrhttp.NewAdapter(
		router,
		mrhttp.WithLogger(opts.Logger),
		mrhttp.WithCaption(restServerCaption),
		mrhttp.WithHostAndPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		mrhttp.WithReadTimeout(srvOpts.ReadTimeout),
		mrhttp.WithWriteTimeout(srvOpts.WriteTimeout),
		mrhttp.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	), nil
}
