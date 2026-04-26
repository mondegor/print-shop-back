package service

import (
	"fmt"

	"github.com/mondegor/go-components/mrmailer/entity"
	"github.com/mondegor/go-components/mrmailer/sendmessage"
	"github.com/mondegor/go-components/mrmailer/sendmessage/adapter"
	"github.com/mondegor/go-components/mrmailer/sendmessage/provider"
	"github.com/mondegor/go-components/mrmailer/service/produce"
	"github.com/mondegor/go-components/wire/mrmailer/processor"
	"github.com/mondegor/go-components/wire/mrmailer/producer"
	"github.com/mondegor/go-components/wire/mrmailer/scheduler"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrclient/mail"
	"github.com/mondegor/go-webcore/mrclient/telegram"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/consume"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceMailerTableName      = "printshop_global.mailer_messages"
	serviceMailerQueueTableName = "printshop_global.mailer_queue"
	serviceMailerPrimaryKey     = "message_id"
)

// InitMailerAPI - создаёт отправителя персонализированных уведомлений получателям.
func InitMailerAPI(opts app.Options) *produce.MessageProducer {
	mrlog.Info(opts.Logger, "Create and init mailer sender API")

	return producer.InitService(
		opts.PostgresConnManager,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceMailerTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceMailerQueueTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		produce.WithRetryAttempts(int16(opts.Cfg.TaskScheduleMailer.SendRetryAttempts)),
		produce.WithDelayCorrection(opts.Cfg.TaskScheduleMailer.SendDelayCorrection),
	)
}

// InitMailerProcessorService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitMailerProcessorService(opts app.Options) (*consume.MessageProcessor[entity.Message], error) {
	mrlog.Info(opts.Logger, "Create and init mail processor service")

	mailSender := sendmessage.NewNopSender()
	messengerSender := sendmessage.NewNopSender()

	mrlog.Info(opts.Logger, "opts.Cfg.MailDefaultFrom", opts.Cfg.MailDefaultFrom)

	if opts.Cfg.MailSmtpHost != "" {
		mrlog.Info(opts.Logger, "Create and init mail client", "host", opts.Cfg.MailSmtpHost, "port", opts.Cfg.MailSmtpPort)

		sender, err := adapter.NewMailSender(
			mail.NewSMTPClient(
				opts.Cfg.MailSmtpHost,
				opts.Cfg.MailSmtpPort,
				opts.Cfg.MailSmtpUserName,
				opts.Cfg.MailSmtpPassword,
				opts.Tracer,
			),
			opts.Cfg.MailDefaultFrom,
		)
		if err != nil {
			return nil, fmt.Errorf("mail.New(): %w", err)
		}

		mailSender = sender
	}

	if opts.Cfg.TelegramChannelName != "" {
		mrlog.Info(opts.Logger, "Create and init telegram bot", "name", opts.Cfg.TelegramChannelName)

		sender, err := telegram.NewBotClient(opts.Cfg.TelegramChannelToken, opts.Tracer)
		if err != nil {
			return nil, fmt.Errorf("telegrambot.NewMessageClient(): %w", err)
		}

		messengerSender = adapter.NewMessengerSender(sender)
	}

	return processor.InitService(
		opts.PostgresConnManager,
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceMailerTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceMailerQueueTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		processor.WithMessageProcessorOpts(
			consume.WithCaptionPrefix[entity.Message]("Mailer/"),
			consume.WithReadyTimeout[entity.Message](opts.Cfg.TaskScheduleMailer.MessageProcessor.ReadyTimeout),
			consume.WithReadPeriodStrategy[entity.Message](
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleMailer.MessageProcessor.ReadPeriod,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			consume.WithConsumerTimeout[entity.Message](
				opts.Cfg.TaskScheduleMailer.MessageProcessor.ConsumerReadTimeout,
				opts.Cfg.TaskScheduleMailer.MessageProcessor.ConsumerWriteTimeout,
			),
			consume.WithHandlerTimeout[entity.Message](opts.Cfg.TaskScheduleMailer.MessageProcessor.HandlerTimeout),
			consume.WithQueueSize[entity.Message](int(opts.Cfg.TaskScheduleMailer.MessageProcessor.QueueSize)),
			consume.WithWorkersCount[entity.Message](int(opts.Cfg.TaskScheduleMailer.MessageProcessor.WorkersCount)),
			consume.WithSignalExecuteHandler[entity.Message](
				opts.PostgresNotificationService.MustFind(opts.Cfg.TaskScheduleMailer.MessageProcessor.NotificationChannel),
			),
		),
		processor.WithSenderProviderOpts(
			provider.WithTracer(opts.Tracer),
			provider.WithClientMail(mailSender),
			provider.WithClientMessenger(messengerSender),
		),
	), nil
}

// InitMailerSchedulerService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitMailerSchedulerService(opts app.Options) *schedule.TaskScheduler {
	mrlog.Info(opts.Logger, "Create and init mail scheduler service")

	return scheduler.InitService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceMailerTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceMailerQueueTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		scheduler.WithCaptionPrefix("Mailer/"),
		scheduler.WithChangeBatchSize(int(opts.Cfg.TaskScheduleMailer.ChangeQueueBatchSize)),
		scheduler.WithChangeRetryTimeout(opts.Cfg.TaskScheduleMailer.ChangeRetryTimeout),
		scheduler.WithChangeRetryDelayed(opts.Cfg.TaskScheduleMailer.ChangeRetryDelayed),
		scheduler.WithCleanBatchSize(int(opts.Cfg.TaskScheduleMailer.CleanQueueBatchSize)),
		scheduler.WithTaskChangeFromToRetryOpts(
			task.WithCaptionPrefix("Mailer/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleMailer.ChangeFromToRetry.Period,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskScheduleMailer.ChangeFromToRetry.Timeout),
		),
		scheduler.WithTaskCleanMessagesOpts(
			task.WithCaptionPrefix("Mailer/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrworker.NewQuadQuickStartStrategy(
					opts.Cfg.TaskScheduleMailer.CleanQueue.Period,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskScheduleMailer.CleanQueue.Timeout),
		),
	)
}
