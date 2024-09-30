package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrhttp"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	internalServerCaption = "HttpInternalServer"
)

// NewInternalServer - создаёт объект mrserver.ServerAdapter.
func NewInternalServer(ctx context.Context, opts app.Options) (*mrhttp.Adapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", internalServerCaption)

	srvOpts := opts.Cfg.Servers.InternalServer

	return mrhttp.NewAdapter(
		ctx,
		opts.InternalRouter,
		mrhttp.WithCaption(internalServerCaption),
		mrhttp.WithHostAndPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		mrhttp.WithReadTimeout(srvOpts.ReadTimeout),
		mrhttp.WithWriteTimeout(srvOpts.WriteTimeout),
		mrhttp.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	), nil
}
