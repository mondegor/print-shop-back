package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mondegor/go-core/mrprocess"
	"github.com/mondegor/go-core/mrprocess/onstartup"
	"github.com/mondegor/go-core/mrprocess/signal"
	"github.com/mondegor/go-core/mrrun"
	"github.com/oklog/run"

	"print-shop-back/cmd/factory"
	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
)

// go get -u github.com/oklog/run

func main() {
	if err := runApp(os.Args, os.Stdout); err != nil {
		if errors.Is(err, config.ErrParseArgsHelp) {
			os.Exit(0)
		}

		log.Fatal(err.Error())
	}
}

func runApp(args []string, stdout io.Writer) error {
	opts, err := factory.InitApp(args, stdout)
	if err != nil {
		return err
	}

	defer func() {
		// close opened shared resources when app shutdown (db, redis, etc...)
		opts.OpenedResources.Close()
		log.Info(opts.Logger, "The application has been stopped")
	}()

	ctx := context.Background()

	appRunner := mrrun.NewAppRunner(&run.Group{}, opts.Logger, opts.TraceManager)
	interceptor := signal.NewInterceptor(opts.Logger)
	lastStarting := appRunner.AddFirstProcess(ctx, interceptor)

	// init services
	{
		for _, service := range opts.TaskSchedulerServices {
			lastStarting = appRunner.AddNextProcess(ctx, service, lastStarting)
		}

		lastStarting = appRunner.AddNextProcess(ctx, opts.PostgresNotificationService, lastStarting)
		lastStarting = appRunner.AddNextProcess(ctx, opts.UserStatRequestCollectorService, lastStarting)
		lastStarting = appRunner.AddNextProcess(ctx, opts.SecureOperationLogCollectorService, lastStarting)
		lastStarting = appRunner.AddNextProcess(ctx, opts.MailProcessorService, lastStarting)
		lastStarting = appRunner.AddNextProcess(ctx, opts.NoticeProcessorService, lastStarting)
	}

	// init app servers
	{
		lastStarting = appRunner.AddNextProcess(ctx, opts.HttpMonitoringServer, lastStarting)
		lastStarting = appRunner.AddNextProcess(ctx, opts.HttpServer, lastStarting)
	}

	// the last process in the startup app chain
	{
		onStartupProcess := onstartup.NewProcess(
			mrprocess.JobFunc(
				func(_ context.Context) error {
					opts.AppHealth.StartupCompleted()
					log.Info(opts.Logger, "The application started, waiting for requests. To exit press CTRL+C")

					return nil
				},
			),
			opts.Logger,
			opts.TraceManager,
		)

		appRunner.AddNextProcess(ctx, onStartupProcess, lastStarting)
	}

	if err = appRunner.Run(); err != nil {
		return fmt.Errorf("the application has been stopped with error: %w", err)
	}

	return nil
}
