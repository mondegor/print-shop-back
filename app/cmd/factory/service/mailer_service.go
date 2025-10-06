package service

import (
	"github.com/mondegor/go-components/factory/mrmailer/processor"
	"github.com/mondegor/go-components/factory/mrmailer/producer"
	"github.com/mondegor/go-components/factory/mrmailer/scheduler"
	"github.com/mondegor/go-components/mrmailer"
	"github.com/mondegor/go-components/mrmailer/provider/mail"
	"github.com/mondegor/go-components/mrmailer/provider/messenger"
	"github.com/mondegor/go-components/mrmailer/provider/nop"
	"github.com/mondegor/go-components/mrmailer/usecase/handle"
	"github.com/mondegor/go-components/mrmailer/usecase/produce"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrsender/mail/smtp"
	"github.com/mondegor/go-webcore/mrsender/telegrambot"
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
func InitMailerAPI(opts app.Options) *produce.MessageSender {
	mrlog.Info(opts.Logger, "Create and init mailer sender API")

	return producer.NewSender(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.UseCaseErrorWrapper,
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

	mailProvider := mrmailer.MessageProvider(nop.New(opts.Tracer))
	telegramProvider := mrmailer.MessageProvider(nop.New(opts.Tracer))

	if opts.Cfg.Senders.Mail.SmtpHost != "" {
		mrlog.Info(opts.Logger, "Create and init mail client", "host", opts.Cfg.Senders.Mail.SmtpHost, "port", opts.Cfg.Senders.Mail.SmtpPort)

		provider, err := mail.New(
			smtp.NewMailClient(
				opts.Cfg.Senders.Mail.SmtpHost,
				opts.Cfg.Senders.Mail.SmtpPort,
				opts.Cfg.Senders.Mail.SmtpUserName,
				opts.Cfg.Senders.Mail.SmtpPassword,
				opts.Tracer,
			),
			opts.Tracer,
			opts.Cfg.Senders.Mail.DefaultFrom,
		)
		if err != nil {
			return nil, err
		}

		mailProvider = provider
	}

	if opts.Cfg.Senders.TelegramBot.Token != "" {
		mrlog.Info(opts.Logger, "Create and init telegram bot", "name", opts.Cfg.Senders.TelegramBot.Name)

		telegramBot, err := telegrambot.NewMessageClient(opts.Cfg.Senders.TelegramBot.Token, opts.Tracer)
		if err != nil {
			return nil, err
		}

		telegramProvider = messenger.New(telegramBot, opts.Tracer)
	}

	return processor.NewService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.UseCaseErrorWrapper,
		opts.Logger,
		opts.Tracer,
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
		processor.WithMessageHandlerOpts(
			handle.WithClientEmail(mailProvider),
			handle.WithClientMessenger(telegramProvider),
		),
	), nil
}

// InitMailerSchedulerService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitMailerSchedulerService(opts app.Options) *schedule.TaskScheduler {
	mrlog.Info(opts.Logger, "Create and init mail scheduler service")

	return scheduler.NewService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.UseCaseErrorWrapper,
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
		scheduler.WithChangeLimit(opts.Cfg.TaskSchedule.Mailer.ChangeQueueLimit),
		scheduler.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeRetryTimeout),
		scheduler.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Mailer.ChangeRetryDelayed),
		scheduler.WithCleanLimit(opts.Cfg.TaskSchedule.Mailer.CleanQueueLimit),
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
