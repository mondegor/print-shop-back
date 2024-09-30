package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/oklog/run"

	"github.com/mondegor/print-shop-back/cmd/factory"
)

// go get -u github.com/oklog/run

func main() {
	ctx := context.Background()
	if err := runApp(ctx, os.Args, os.Stdout); err != nil {
		if errors.Is(err, factory.ErrParseArgsHelp) {
			os.Exit(0)
		}

		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func runApp(ctx context.Context, args []string, stdout io.Writer) error {
	ctx, opts, err := factory.InitApp(ctx, args, stdout)
	if err != nil {
		return err
	}

	// close opened shared resources when app shutdown (db, redis, etc...)
	defer mrlib.CallEachFunc(ctx, opts.OpenedResources)

	appRunner := mrrun.NewAppRunner(&run.Group{})
	ctx, chLast := appRunner.AddSignalHandler(ctx)

	// init task scheduler
	{
		scheduler, err := factory.NewTaskScheduler(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewScheduler(): %w", err)
		}

		chLast = appRunner.AddNextProcess(ctx, scheduler, chLast)
	}

	// init app servers
	{
		restServer, err := factory.NewRestServer(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewRestServer(): %w", err)
		}

		chLast = appRunner.AddNextProcess(ctx, restServer, chLast)

		internalServer, err := factory.NewInternalServer(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewInternalServer(): %w", err)
		}

		appRunner.AddNextProcess(ctx, internalServer, chLast)
	}

	// run app and its servers
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Starting the application...")

	onStartup := func() {
		logger.Info().Msg("The application started, waiting for requests. To exit press CTRL+C")
		opts.AppHealth.StartupCompleted()
	}

	if err = appRunner.Run(onStartup); err != nil {
		return fmt.Errorf("the application has been stopped with error: %w", err)
	}

	logger.Info().Msg("The application has been stopped")

	return nil
}
