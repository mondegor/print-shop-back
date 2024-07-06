package main

import (
	"flag"
	"log"

	"github.com/mondegor/print-shop-back/cmd/factory"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/oklog/run"
)

// go get -u github.com/oklog/run

func main() {
	ctx, opts, err := factory.CreateAppEnvironment(parseFlags())
	if err != nil {
		log.Fatal(err)
	}

	logger := mrlog.Ctx(ctx)

	opts, err = factory.InitAppEnvironment(ctx, opts)
	if err != nil {
		logger.Fatal().Err(err).Msg("factory.InitAppEnvironment() error")
	}

	// close opened resources when app shutdown (db, redis, etc...)
	defer mrlib.CallEachFunc(ctx, opts.OpenedResources)

	appRunner := run.Group{}
	appRunner.Add(mrserver.PrepareAppToStart(ctx))

	// init task scheduler
	if scheduler, err := factory.NewTaskScheduler(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("factory.NewScheduler() error")
	} else {
		appRunner.Add(scheduler.PrepareToStart(ctx))
	}

	// init app servers
	if restServer, err := factory.NewRestServer(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("factory.NewRestServer() error")
	} else {
		appRunner.Add(restServer.PrepareToStart(ctx))
	}

	if prometheusServer, err := factory.NewPrometheusServer(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("factory.NewPrometheusServer() error")
	} else {
		appRunner.Add(prometheusServer.PrepareToStart(ctx))
	}

	// run app and its servers
	logger.Info().Msg("Running the application...")

	if err = appRunner.Run(); err != nil {
		logger.Error().Err(err).Msg("The application has been stopped with error")
	} else {
		logger.Info().Msg("The application has been stopped")
	}
}

func parseFlags() (configPath, logLevel string) {
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "Path to application config file")
	flag.StringVar(&logLevel, "log-level", "", "Logging level")
	flag.Parse()

	return configPath, logLevel
}
