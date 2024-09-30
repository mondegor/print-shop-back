package factory

import (
	"context"
	"errors"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker/mrschedule"

	"github.com/mondegor/print-shop-back/internal/app"
)

// NewTaskScheduler - создаёт объект mrschedule.Scheduler.
func NewTaskScheduler(ctx context.Context, opts app.Options) (*mrschedule.Scheduler, error) {
	logger := mrlog.Ctx(ctx)
	logger.Info().Msg("Create and init scheduler with it's tasks")

	if len(opts.SchedulerTasks) == 0 {
		return nil, errors.New("opts.SchedulerTasks is empty")
	}

	for _, task := range opts.SchedulerTasks {
		logger.Debug().Msgf("- registered: task %s, period: %s, timeout: %s", task.Caption(), task.Period(), task.Timeout())
	}

	return mrschedule.NewScheduler(opts.ErrorHandler, opts.SchedulerTasks...), nil
}
