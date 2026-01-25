package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/httpserver"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	internalServerCaption = "HttpInternalServer"
)

// InitInternalServer - создаёт объект mrserver.ServerAdapter.
func InitInternalServer(opts app.Options) *httpserver.Adapter {
	srvOpts := opts.Cfg.Servers.InternalServer

	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", internalServerCaption), "port", srvOpts.Listen.Port)

	return httpserver.New(
		opts.InternalRouter,
		httpserver.WithLogger(opts.Logger),
		httpserver.WithCaption(internalServerCaption),
		httpserver.WithHostPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		httpserver.WithReadTimeout(srvOpts.ReadTimeout),
		httpserver.WithWriteTimeout(srvOpts.WriteTimeout),
		httpserver.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	)
}
