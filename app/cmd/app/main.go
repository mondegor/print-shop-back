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
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/mondegor/go-webcore/mrworker/process/onstartup"
	"github.com/mondegor/go-webcore/mrworker/process/signal"
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

	defer func() {
		// close opened shared resources when app shutdown (db, redis, etc...)
		mrlib.CallEachFunc(ctx, opts.OpenedResources)
		mrlog.Ctx(ctx).Info().Msg("The application has been stopped")
	}()

	appRunner := mrrun.NewAppRunner(&run.Group{})
	ctx, interception := signal.NewInterception(ctx)
	lastStarting := appRunner.AddFirstProcess(ctx, interception)

	// init services
	{
		mailer, tasks := factory.NewMailerService(ctx, opts)
		opts.SchedulerTasks = append(opts.SchedulerTasks, tasks...)
		lastStarting = appRunner.AddNextProcess(ctx, mailer, lastStarting)

		notifier, tasks := factory.NewNotifierService(ctx, opts)
		opts.SchedulerTasks = append(opts.SchedulerTasks, tasks...)
		lastStarting = appRunner.AddNextProcess(ctx, notifier, lastStarting)

		scheduler, err := factory.NewTaskScheduler(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewScheduler(): %w", err)
		}

		lastStarting = appRunner.AddNextProcess(ctx, scheduler, lastStarting)
	}

	// init app servers
	{
		restServer, err := factory.NewRestServer(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewRestServer(): %w", err)
		}

		lastStarting = appRunner.AddNextProcess(ctx, restServer, lastStarting)

		internalServer, err := factory.NewInternalServer(ctx, opts)
		if err != nil {
			return fmt.Errorf("factory.NewInternalServer(): %w", err)
		}

		lastStarting = appRunner.AddNextProcess(ctx, internalServer, lastStarting)
	}

	// the last process in the startup app chain
	{
		onStartupProcess := onstartup.NewProcess(
			mrworker.JobFunc(
				func(ctx context.Context) error {
					opts.AppHealth.StartupCompleted()
					mrlog.Ctx(ctx).Info().Msg("The application started, waiting for requests. To exit press CTRL+C")

					return nil
				},
			),
		)

		appRunner.AddNextProcess(ctx, onStartupProcess, lastStarting)
	}

	if err = appRunner.Run(); err != nil {
		return fmt.Errorf("the application has been stopped with error: %w", err)
	}

	return nil
}
