package factory

import (
	"context"
	"errors"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewTaskScheduler - создаёт объект schedule.TaskScheduler.
func NewTaskScheduler(ctx context.Context, opts app.Options) (*schedule.TaskScheduler, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init scheduler with it's tasks")

	if len(opts.SchedulerTasks) == 0 {
		return nil, errors.New("opts.SchedulerTasks is empty")
	}

	for _, task := range opts.SchedulerTasks {
		logger.Debug().Msgf(
			"- registered: task %s, period: %s, timeout: %s; startup: %t",
			task.Caption(),
			task.Period(),
			task.Timeout(),
			task.Startup(),
		)
	}

	return schedule.NewTaskScheduler(
		opts.ErrorHandler,
		schedule.WithTasks(opts.SchedulerTasks...),
	), nil
}
