package factory

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
)

// RegisterSystemHandlers - регистрация системных обработчиков.
func RegisterSystemHandlers(opts app.Options) error {
	mrlog.Info(opts.Logger, "Init system handlers")

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

	// :TODO: логировать регистрацию обработчиков
	opts.InternalRouter.Handle(
		"/health",
		mrresp.HandlerGetHealth(
			mrrun.PrepareProbesForCheck(opts.Logger, probes...),
		),
	)

	probesFunc := mrrun.PrepareProbes(opts.Logger, probes...)

	systemInfoFunc, err := mrresp.HandlerGetSystemInfoAsJSON(
		opts.Logger,
		mrresp.SystemInfoConfig{
			Caption:     opts.Cfg.App.Name,
			Version:     opts.Cfg.App.Version,
			Environment: opts.Cfg.App.Environment,
			IsDebug:     opts.Cfg.Debugging.Debug,
			LogLevel:    opts.Cfg.Log.Level,
			StartedAt:   opts.Cfg.App.StartedAt,
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

	opts.InternalRouter.Handle("/system-info", systemInfoFunc)

	return nil
}
