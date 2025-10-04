package factory

import (
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/mondegor/print-shop-back/internal/app"
)

// go get -u github.com/prometheus/client_golang

// InitPrometheus - создаёт объект mrinit.Prometheus.
func InitPrometheus(opts app.Options) *mrinit.Prometheus {
	prom := mrinit.NewPrometheus()

	prom.Add(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	opts.InternalRouter.Handle(
		"/metrics",
		promhttp.HandlerFor(
			prom.Registry(),
			promhttp.HandlerOpts{
				Registry: prom.Registry(),
			},
		),
	)

	return prom
}
