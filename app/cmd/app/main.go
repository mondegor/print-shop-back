package main

import (
	"flag"
	"print-shop-back/internal/factory"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrserver"
	"github.com/oklog/run"
)

// go get -u github.com/oklog/run

var (
	configPath string
	logLevel   string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "Path to application config file")
	flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
	flag.Parse()

	ctx, opts := factory.CreateAppEnvironment(configPath, logLevel) // success or fatal
	logger := mrlog.Ctx(ctx)

	opts, err := factory.InitAppEnvironment(ctx, opts)

	if err != nil {
		logger.Fatal().Err(err).Msg("factory.InitAppEnvironment() error")
	}

	// close opened resources when app shutdown (db, redis, etc...)
	defer mrlib.CallEachFunc(ctx, opts.OpenedResources)

	appRunner := &run.Group{}
	appRunner.Add(mrserver.PrepareAppToStart(ctx))

	// init app servers
	if restServer, err := factory.NewRestServer(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("factory.NewRestServer() error")
	} else {
		appRunner.Add(restServer.PrepareToStart(ctx))
	}

	// run app and its servers
	logger.Info().Msg("Running the application...")

	if err = appRunner.Run(); err != nil {
		logger.Error().Err(err).Msg("The application has been stopped with error")
	} else {
		logger.Info().Msg("The application has been stopped")
	}
}
