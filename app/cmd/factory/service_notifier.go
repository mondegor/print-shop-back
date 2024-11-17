package factory

import (
	"context"

	"github.com/mondegor/go-components/factory/mrnotifier"
	"github.com/mondegor/go-components/mrnotifier/notifier/component/produce"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/consume"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceNotifierTableName             = "printshop_global.notifier_notices"
	serviceNotifierQueueTableName        = "printshop_global.notifier_queue"
	serviceNotifierPrimaryKey            = "notice_id"
	serviceNotifierTemplatesTableName    = "printshop_global.notifier_templates"
	serviceNotifierTemplateVarsTableName = "printshop_global.notifier_template_vars"
)

// NewNotifierAPI - создаёт отправителя сообщений получателям.
func NewNotifierAPI(ctx context.Context, opts app.Options) *produce.NoticeSender {
	mrlog.Ctx(ctx).Info().Msg("Create and init notifier sender API")

	return mrnotifier.NewComponentSender(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceNotifierTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceNotifierQueueTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		opts.EventEmitter,
		produce.WithRetryAttempts(opts.Cfg.TaskSchedule.Notifier.SendRetryAttempts),
	)
}

// NewNotifierService - создаёт сервис для обработки уведомлений и связанных с ним задачи.
func NewNotifierService(ctx context.Context, opts app.Options) (mrrun.Process, []mrworker.Task) {
	mrlog.Ctx(ctx).Info().Msg("Create and init notifier service")

	return mrnotifier.NewComponentService(
		opts.PostgresConnManager,
		opts.MailerAPI,
		opts.EventEmitter,
		opts.ErrorHandler,
		mrsql.DBTableInfo{
			Name:       serviceNotifierTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceNotifierQueueTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		serviceNotifierTemplatesTableName,
		serviceNotifierTemplateVarsTableName,
		mrnotifier.WithSendProcessorOpts(
			consume.WithCaption(opts.Cfg.TaskSchedule.Notifier.SendProcessor.Caption),
			consume.WithReadyTimeout(opts.Cfg.TaskSchedule.Notifier.SendProcessor.ReadyTimeout),
			consume.WithReadPeriod(opts.Cfg.TaskSchedule.Notifier.SendProcessor.ReadPeriod),
			consume.WithBusyReadPeriod(opts.Cfg.TaskSchedule.Notifier.SendProcessor.BusyReadPeriod),
			consume.WithCancelReadTimeout(opts.Cfg.TaskSchedule.Notifier.SendProcessor.CancelReadTimeout),
			consume.WithHandlerTimeout(opts.Cfg.TaskSchedule.Notifier.SendProcessor.HandlerTimeout),
			consume.WithQueueSize(opts.Cfg.TaskSchedule.Notifier.SendProcessor.QueueSize),
			consume.WithWorkersCount(opts.Cfg.TaskSchedule.Notifier.SendProcessor.WorkersCount),
		),
		mrnotifier.WithTaskChangeFromToRetryOpts(
			task.WithCaption(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Caption),
			task.WithStartup(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Startup),
			task.WithPeriod(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Timeout),
		),
		mrnotifier.WithTaskCleanNoticesOpts(
			task.WithCaption(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Caption),
			task.WithStartup(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Startup),
			task.WithPeriod(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Timeout),
		),
		mrnotifier.WithChangeLimit(opts.Cfg.TaskSchedule.Notifier.ChangeQueueLimit),
		mrnotifier.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeRetryTimeout),
		mrnotifier.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Mailer.ChangeRetryDelayed),
		mrnotifier.WithCleanLimit(opts.Cfg.TaskSchedule.Notifier.CleanQueueLimit),
		mrnotifier.WithDefaultLang(opts.Cfg.Translation.LangCodes[0]),
	)
}
