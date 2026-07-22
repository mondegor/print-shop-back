package factory

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mondegor/go-core/mrrun"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
)

// RegisterSystemHandlers - регистрация системных обработчиков.
func RegisterSystemHandlers(opts app.Options) error {
	log.Info(opts.Logger, "Init system handlers")

	probes := []mrrun.ProbeChecker{
		mrrun.NewHealthProbe(
			opts.Logger,
			"App",
			mrrun.WithAppReadyProbe(opts.AppHealth),
			time.Microsecond,
		),
		mrrun.NewHealthProbe(
			opts.Logger,
			"Postgres",
			opts.PostgresConnManager.ConnAdapter().Ping,
			0, // timeout by default
		),
		mrrun.NewHealthProbe(
			opts.Logger,
			"Redis",
			opts.RedisAdapter.Ping,
			0, // timeout by default
		),
	}

	// TODO: логировать регистрацию обработчиков
	opts.MonitoringRouter.Handle(
		"/health",
		mrresp.HandlerGetHealth(
			mrrun.PrepareProbesForCheck(opts.Logger, probes...),
		),
	)

	probesFunc := mrrun.PrepareProbes(opts.Logger, probes...)

	systemInfoFunc, err := mrresp.HandlerGetSystemInfoAsJSON(
		opts.Logger,
		mrresp.SystemInfoConfig{
			Caption:     opts.Cfg.AppName,
			Version:     opts.Cfg.AppVersion,
			Environment: opts.Cfg.Environment,
			IsDebug:     opts.Cfg.DebugIsEnabled,
			LogLevel:    opts.Cfg.LogLevel,
			StartedAt:   opts.Cfg.StartedAt,
			ProcessesFunc: func(ctx context.Context) []mrresp.SystemInfoProcess {
				finishedProbes := probesFunc(ctx)
				processes := make([]mrresp.SystemInfoProcess, 0, len(finishedProbes))

				for _, probe := range finishedProbes {
					processes = append(
						processes,
						mrresp.SystemInfoProcess{
							Caption: probe.Caption,
							Status:  strconv.Itoa(probe.Status) + " " + http.StatusText(probe.Status),
						},
					)
				}

				return processes
			},
		},
	)
	if err != nil {
		return err
	}

	opts.MonitoringRouter.Handle("/v1/system-info", systemInfoFunc)

	return nil
}
