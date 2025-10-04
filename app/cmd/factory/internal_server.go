package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrhttp"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	internalServerCaption = "HttpInternalServer"
)

// InitInternalServer - создаёт объект mrserver.ServerAdapter.
func InitInternalServer(opts app.Options) *mrhttp.Adapter {
	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", internalServerCaption))

	srvOpts := opts.Cfg.Servers.InternalServer

	return mrhttp.NewAdapter(
		opts.InternalRouter,
		mrhttp.WithLogger(opts.Logger),
		mrhttp.WithCaption(internalServerCaption),
		mrhttp.WithHostAndPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		mrhttp.WithReadTimeout(srvOpts.ReadTimeout),
		mrhttp.WithWriteTimeout(srvOpts.WriteTimeout),
		mrhttp.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	)
}
