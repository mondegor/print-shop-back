package factory

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrserver/mrresp"

	"github.com/mondegor/print-shop-back/internal/app"
)

// RegisterSystemHandlers - регистрация системных обработчиков.
func RegisterSystemHandlers(ctx context.Context, opts app.Options) error {
	mrlog.Ctx(ctx).Info().Msgf("Init system handlers")

	probes := []mrrun.ProbeChecker{
		mrrun.NewHealthProbe(
			"App",
			mrrun.WithAppReadyProbe(opts.AppHealth),
			time.Microsecond,
		),
		mrrun.NewHealthProbe(
			"Postgres",
			opts.PostgresConnManager.ConnAdapter().Ping,
			0, // timeout by default
		),
		mrrun.NewHealthProbe(
			"Redis",
			opts.RedisAdapter.Ping,
			0, // timeout by default
		),
	}

	// :TODO: логировать регистрацию обработчиков
	opts.InternalRouter.Handle(
		"/health",
		mrresp.HandlerGetHealth(
			mrrun.PrepareProbesForCheck(probes...),
		),
	)

	probesFunc := mrrun.PrepareProbes(probes...)

	systemInfoFunc, err := mrresp.HandlerGetSystemInfoAsJSON(
		mrresp.SystemInfoConfig{
			Name:        opts.Cfg.App.Name,
			Version:     opts.Cfg.App.Version,
			Environment: opts.Cfg.App.Environment,
			IsDebug:     opts.Cfg.Debugging.Debug,
			LogLevel:    mrlog.Ctx(ctx).Level(),
			StartedAt:   opts.Cfg.App.StartedAt,
			Processes: func(ctx context.Context) map[string]string {
				finishedProbes := probesFunc(ctx)
				processes := make(map[string]string, len(finishedProbes))

				for _, probe := range finishedProbes {
					processes[probe.Caption] = strconv.Itoa(probe.Status) + " " + http.StatusText(probe.Status)
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
