package service

import (
	"github.com/mondegor/go-components/mrnotifier/notifier/entity"
	"github.com/mondegor/go-components/mrnotifier/notifier/service/produce"
	"github.com/mondegor/go-components/wire/mrmailer"
	"github.com/mondegor/go-components/wire/mrnotifier/processor"
	"github.com/mondegor/go-components/wire/mrnotifier/producer"
	"github.com/mondegor/go-components/wire/mrnotifier/scheduler"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
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
func InitNotifierAPI(opts app.Options) *produce.NoteProducer {
	mrlog.Info(opts.Logger, "Create and init notifier sender API")

	return producer.InitService(
		opts.PostgresConnManager,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceNotifierTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceNotifierQueueTableName,
			PrimaryKey: serviceNotifierPrimaryKey,
		},
		produce.WithRetryAttempts(int16(opts.Cfg.TaskSchedule.Notifier.SendRetryAttempts)),
	)
}

// InitNotifierProcessorService - создаёт сервис для обработки уведомлений и связанных с ним задачи.
func InitNotifierProcessorService(opts app.Options) *consume.MessageProcessor[entity.Note] {
	mrlog.Info(opts.Logger, "Create and init notice processor service")

	return processor.InitService(
		opts.PostgresConnManager,
		mrmailer.NoticeToMessageAdapterFunc(opts.MailerAPI.Send),
		opts.ErrorHandler,
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
			consume.WithCaptionPrefix[entity.Note]("Notifier/"),
			consume.WithReadyTimeout[entity.Note](opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ReadyTimeout),
			consume.WithReadPeriodStrategy[entity.Note](
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ReadPeriod,
					opts.Cfg.TaskSchedule.Settings.DefaultPeriodRatio,
				),
			),
			consume.WithConsumerTimeout[entity.Note](
				opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ConsumerReadTimeout,
				opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.ConsumerWriteTimeout,
			),
			consume.WithHandlerTimeout[entity.Note](opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.HandlerTimeout),
			consume.WithQueueSize[entity.Note](int(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.QueueSize)),
			consume.WithWorkersCount[entity.Note](int(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.WorkersCount)),
			consume.WithSignalExecuteHandler[entity.Note](
				opts.PostgresNotificationService.ReceiverChannels.MustFind(opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.NotificationChannel),
			),
		),
	)
}

// InitNotifierSchedulerService - создаёт сервис для обработки уведомлений и связанных с ним задачи.
func InitNotifierSchedulerService(opts app.Options) *schedule.TaskScheduler {
	mrlog.Info(opts.Logger, "Create and init notice scheduler service")

	return scheduler.InitService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
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
		scheduler.WithChangeBatchSize(int(opts.Cfg.TaskSchedule.Notifier.ChangeQueueBatchSize)),
		scheduler.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Notifier.ChangeRetryTimeout),
		scheduler.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Notifier.ChangeRetryDelayed),
		scheduler.WithCleanBatchSize(int(opts.Cfg.TaskSchedule.Notifier.CleanQueueBatchSize)),
		scheduler.WithTaskChangeFromToRetryOpts(
			task.WithCaptionPrefix("Notifier/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Period,
					opts.Cfg.TaskSchedule.Settings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.ChangeFromToRetry.Timeout),
		),
		scheduler.WithTaskCleanNoticesOpts(
			task.WithCaptionPrefix("Notifier/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrworker.NewQuadQuickStartStrategy(
					opts.Cfg.TaskSchedule.Notifier.CleanQueue.Period,
					opts.Cfg.TaskSchedule.Settings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskSchedule.Notifier.CleanQueue.Timeout),
		),
	)
}
