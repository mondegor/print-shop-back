package factory

import (
	"fmt"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrserver/httpserver"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	monitoringServerCaption = "HttpMonitoringServer"
)

// InitMonitoringServer - создаёт объект mrserver.ServerAdapter.
func InitMonitoringServer(opts app.Options) *httpserver.Adapter {
	mrlog.Info(opts.Logger, fmt.Sprintf("Create and init '%s'", monitoringServerCaption), "port", opts.Cfg.MonitoringServerPort)

	return httpserver.New(
		opts.MonitoringRouter,
		httpserver.WithLogger(opts.Logger),
		httpserver.WithCaption(monitoringServerCaption),
		httpserver.WithHostPort(opts.Cfg.MonitoringServerBindIP, opts.Cfg.MonitoringServerPort),
		httpserver.WithReadTimeout(opts.Cfg.MonitoringServerReadTimeout),
		httpserver.WithWriteTimeout(opts.Cfg.MonitoringServerWriteTimeout),
		httpserver.WithShutdownTimeout(opts.Cfg.MonitoringServerShutdownTimeout),
	)
}
