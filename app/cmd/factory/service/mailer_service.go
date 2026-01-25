package service

import (
	"fmt"

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
		produce.WithRetryAttempts(opts.Cfg.TaskSchedule.Mailer.SendRetryAttempts),
		produce.WithDelayCorrection(opts.Cfg.TaskSchedule.Mailer.SendDelayCorrection),
	)
}

// InitMailerProcessorService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitMailerProcessorService(opts app.Options) (*consume.MessageProcessor, error) {
	mrlog.Info(opts.Logger, "Create and init mail processor service")

	mailSender := sendmessage.NewNopSender()
	messengerSender := sendmessage.NewNopSender()

	mrlog.Info(opts.Logger, "opts.Cfg.Senders.Mail.DefaultFrom", opts.Cfg.Senders.Mail.DefaultFrom)

	if opts.Cfg.Senders.Mail.SmtpHost != "" {
		mrlog.Info(opts.Logger, "Create and init mail client", "host", opts.Cfg.Senders.Mail.SmtpHost, "port", opts.Cfg.Senders.Mail.SmtpPort)

		sender, err := adapter.NewMailSender(
			mail.NewSMTPClient(
				opts.Cfg.Senders.Mail.SmtpHost,
				opts.Cfg.Senders.Mail.SmtpPort,
				opts.Cfg.Senders.Mail.SmtpUserName,
				opts.Cfg.Senders.Mail.SmtpPassword,
				opts.Tracer,
			),
			opts.Cfg.Senders.Mail.DefaultFrom,
		)
		if err != nil {
			return nil, fmt.Errorf("mail.New(): %w", err)
		}

		mailSender = sender
	}

	if opts.Cfg.Senders.TelegramBot.Token != "" {
		mrlog.Info(opts.Logger, "Create and init telegram bot", "name", opts.Cfg.Senders.TelegramBot.Name)

		sender, err := telegram.NewBotClient(opts.Cfg.Senders.TelegramBot.Token, opts.Tracer)
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
			consume.WithCaptionPrefix("Mailer/"),
			consume.WithReadyTimeout(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.ReadyTimeout),
			consume.WithReadPeriod(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.ReadPeriod),
			consume.WithConsumerTimeout(
				opts.Cfg.TaskSchedule.Mailer.MessageProcessor.ConsumerReadTimeout,
				opts.Cfg.TaskSchedule.Mailer.MessageProcessor.ConsumerWriteTimeout,
			),
			consume.WithHandlerTimeout(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.HandlerTimeout),
			consume.WithQueueSize(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.QueueSize),
			consume.WithWorkersCount(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.WorkersCount),
			consume.WithSignalExecuteHandler(
				opts.PostgresNotificationService.ReceiverChannels.MustFind(opts.Cfg.TaskSchedule.Mailer.MessageProcessor.NotificationChannel),
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
		scheduler.WithChangeBatchSize(opts.Cfg.TaskSchedule.Mailer.ChangeQueueBatchSize),
		scheduler.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeRetryTimeout),
		scheduler.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Mailer.ChangeRetryDelayed),
		scheduler.WithCleanBatchSize(opts.Cfg.TaskSchedule.Mailer.CleanQueueBatchSize),
		scheduler.WithTaskChangeFromToRetryOpts(
			task.WithCaptionPrefix("Mailer/"),
			task.WithStartup(false),
			task.WithPeriod(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Timeout),
		),
		scheduler.WithTaskCleanMessagesOpts(
			task.WithCaptionPrefix("Mailer/"),
			task.WithStartup(false),
			task.WithPeriod(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Timeout),
		),
	)
}
