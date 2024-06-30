package factory

import (
	"context"
	"net/http"

	"github.com/mondegor/print-shop-back/internal/app"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	prometheusServerCaption = "Prometheus"
)

// NewPrometheusServer - создаёт объект mrserver.ServerAdapter.
func NewPrometheusServer(ctx context.Context, opts app.Options) (*mrserver.ServerAdapter, error) {
	mrlog.Ctx(ctx).Info().Msgf("Create and init '%s'", prometheusServerCaption)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.HandlerFor(opts.Prometheus, promhttp.HandlerOpts{Registry: opts.Prometheus}))

	srvOpts := opts.Cfg.Servers.PrometheusServer

	return mrserver.NewServerAdapter(
		ctx,
		mrserver.ServerOptions{
			Caption:         prometheusServerCaption,
			Handler:         router,
			ReadTimeout:     srvOpts.ReadTimeout,
			WriteTimeout:    srvOpts.WriteTimeout,
			ShutdownTimeout: srvOpts.ShutdownTimeout,
			Listen: mrserver.ListenOptions{
				BindIP: srvOpts.Listen.BindIP,
				Port:   srvOpts.Listen.Port,
			},
		},
	), nil
}
