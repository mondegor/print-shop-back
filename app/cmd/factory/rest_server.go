package factory

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver/mrhttp"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	restServerCaption = "RestServer"
)

// NewRestServer - создаёт объект mrserver.ServerAdapter.
func NewRestServer(ctx context.Context, opts app.Options) (*mrhttp.Adapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", restServerCaption)

	router, err := NewRestRouterWithRegisterHandlers(ctx, opts)
	if err != nil {
		return nil, err
	}

	srvOpts := opts.Cfg.Servers.RestServer

	return mrhttp.NewAdapter(
		ctx,
		router,
		mrhttp.WithCaption(restServerCaption),
		mrhttp.WithHostAndPort(srvOpts.Listen.BindIP, srvOpts.Listen.Port),
		mrhttp.WithReadTimeout(srvOpts.ReadTimeout),
		mrhttp.WithWriteTimeout(srvOpts.WriteTimeout),
		mrhttp.WithShutdownTimeout(srvOpts.ShutdownTimeout),
	), nil
}
