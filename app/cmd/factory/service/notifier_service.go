package service

import (
	"github.com/mondegor/go-components/factory/mrnotifier/processor"
	"github.com/mondegor/go-components/factory/mrnotifier/producer"
	"github.com/mondegor/go-components/factory/mrnotifier/scheduler"
	"github.com/mondegor/go-components/mrnotifier/notifier/usecase/produce"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/consume"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceNotifierTableName             = "printshop_global.notifier_notices"
	serviceNotifierQueueTableName        = "printshop_global.notifier_queue"
	serviceNotifierPrimaryKey            = "notice_id"
	serviceNotifierTemplatesTableName    = "printshop_global.notifier_templates"
	serviceNotifierTemplateVarsTableName = "printshop_global.notifier_template_vars"
)

// InitNotifierAPI - создаёт отправителя сообщений получателям.
func InitNotifierAPI(opts app.Options) *produce.NoticeSender {
	mrlog.Info(opts.Logger, "Create and init notifier sender API")

	return producer.NewSender(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceNotifierTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceNotifierQueueTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		produce.WithRetryAttempts(opts.Cfg.TaskSchedule.Notifier.SendRetryAttempts),
	)
}

// InitNotifierProcessorService - создаёт сервис для обработки уведомлений и связанных с ним задачи.
func InitNotifierProcessorService(opts app.Options) *consume.MessageProcessor {
	mrlog.Info(opts.Logger, "Create and init notice processor service")

	return processor.NewService(
		opts.PostgresConnManager,
		opts.MailerAPI,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.UsecaseErrorWrapper,
		opts.StorageErrorWrapper,
		opts.Logger,
		opts.TraceManager,
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
		processor.WithDefaultLang(opts.Cfg.Localization.Languages[0]),
		processor.WithNoticeProcessorOpts(
			consume.WithCaptionPrefix("Notifier/"),
			consume.WithReadyTimeout(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ReadyTimeout),
			consume.WithReadPeriod(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ReadPeriod),
			consume.WithConsumerTimeout(
				opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ConsumerReadTimeout,
				opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ConsumerWriteTimeout,
			),
			consume.WithHandlerTimeout(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.HandlerTimeout),
			consume.WithQueueSize(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.QueueSize),
			consume.WithWorkersCount(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.WorkersCount),
			consume.WithSignalExecuteHandler(
				opts.PostgresNotificationService.ReceiverChannels.MustFind(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.NotificationChannel),
			),
		),
	)
}

// InitNotifierSchedulerService - создаёт сервис для обработки уведомлений и связанных с ним задачи.
func InitNotifierSchedulerService(opts app.Options) *schedule.TaskScheduler {
	mrlog.Info(opts.Logger, "Create and init notice scheduler service")

	return scheduler.NewService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.UsecaseErrorWrapper,
		opts.Logger,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceNotifierTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceNotifierQueueTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		scheduler.WithCaptionPrefix("Notifier/"),
		scheduler.WithChangeLimit(opts.Cfg.TaskSchedule.Notifier.ChangeQueueLimit),
		scheduler.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Notifier.ChangeRetryTimeout),
		scheduler.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Notifier.ChangeRetryDelayed),
		scheduler.WithCleanLimit(opts.Cfg.TaskSchedule.Notifier.CleanQueueLimit),
		scheduler.WithTaskChangeFromToRetryOpts(
			task.WithCaptionPrefix("Notifier/"),
			task.WithStartup(false),
			task.WithPeriod(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Timeout),
		),
		scheduler.WithTaskCleanNoticesOpts(
			task.WithCaptionPrefix("Notifier/"),
			task.WithStartup(false),
			task.WithPeriod(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Timeout),
		),
	)
}
