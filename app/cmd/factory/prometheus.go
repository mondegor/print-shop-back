package factory

import (
	"github.com/mondegor/go-core/mrstorage"
	"github.com/mondegor/go-webcore/mrcore/mrinit"
	"github.com/mondegor/go-webcore/mrstorage/mrprometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
)

// InitPrometheus - создаёт объект mrinit.Prometheus.
func InitPrometheus(opts app.Options) *mrinit.Prometheus {
	prom := mrinit.NewPrometheus()

	prom.Add(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	opts.MonitoringRouter.Handle(
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

// InitPrometheusStatPostgres - инициализирует сбор информации о Postgres.
func InitPrometheusStatPostgres(opts app.Options) error {
	if opts.Prometheus == nil || opts.PostgresConnManager == nil {
		log.Warn(opts.Logger, "Collector DB Stat is disabled")

		return nil
	}

	connCli, err := opts.PostgresConnManager.ConnAdapter().Cli()
	if err != nil {
		return err
	}

	opts.Prometheus.Add(
		mrprometheus.NewDBCollector(
			"pgx",
			func() mrstorage.DBStatProvider {
				return connCli.Stat()
			},
			map[string]string{
				"db_name": opts.Cfg.DBDatabase,
			},
		),
	)

	return nil
}
