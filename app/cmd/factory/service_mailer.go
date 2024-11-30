package factory

import (
	"context"

	"github.com/mondegor/go-components/factory/mrmailer"
	mailer "github.com/mondegor/go-components/mrmailer"
	"github.com/mondegor/go-components/mrmailer/component/handle"
	"github.com/mondegor/go-components/mrmailer/component/produce"
	"github.com/mondegor/go-components/mrmailer/provider/mail"
	"github.com/mondegor/go-components/mrmailer/provider/messenger"
	"github.com/mondegor/go-components/mrmailer/provider/nop"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrrun"
	"github.com/mondegor/go-webcore/mrsender/mail/smtp"
	"github.com/mondegor/go-webcore/mrsender/telegrambot"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/consume"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceMailerTableName      = "printshop_global.mailer_messages"
	serviceMailerQueueTableName = "printshop_global.mailer_queue"
	serviceMailerPrimaryKey     = "message_id"
)

// NewMailerAPI - создаёт отправителя персонализированных уведомлений получателям.
func NewMailerAPI(ctx context.Context, opts app.Options) *produce.MessageSender {
	mrlog.Ctx(ctx).Info().Msg("Create and init mailer sender API")

	return mrmailer.NewComponentSender(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceMailerTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceMailerQueueTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		opts.EventEmitter,
		produce.WithRetryAttempts(opts.Cfg.TaskSchedule.Mailer.SendRetryAttempts),
		produce.WithDelayCorrection(opts.Cfg.TaskSchedule.Mailer.SendDelayCorrection),
	)
}

// NewMailerService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func NewMailerService(ctx context.Context, opts app.Options) (mrrun.Process, []mrworker.Task, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init mailer service")

	mailProvider := mailer.MessageProvider(nop.New())
	telegramProvider := mailer.MessageProvider(nop.New())

	if opts.Cfg.Senders.Mail.SmtpHost != "" {
		provider, err := mail.New(
			smtp.NewMailClient(
				opts.Cfg.Senders.Mail.SmtpHost,
				opts.Cfg.Senders.Mail.SmtpPort,
				opts.Cfg.Senders.Mail.SmtpUserName,
				opts.Cfg.Senders.Mail.SmtpPassword,
			),
			opts.Cfg.Senders.Mail.DefaultFrom,
		)
		if err != nil {
			return nil, nil, err
		}

		mailProvider = provider
	}

	if opts.Cfg.Senders.TelegramBot.Token != "" {
		telegramBot, err := telegrambot.NewMessageClient(opts.Cfg.Senders.TelegramBot.Token)
		if err != nil {
			return nil, nil, err
		}

		telegramProvider = messenger.New(telegramBot)
	}

	process, tasks := mrmailer.NewComponentService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		mrsql.DBTableInfo{
			Name:       serviceMailerTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceMailerQueueTableName,
			PrimaryKey: serviceMailerPrimaryKey,
		},
		mrmailer.WithSendProcessorOpts(
			consume.WithCaption(opts.Cfg.TaskSchedule.Mailer.SendProcessor.Caption),
			consume.WithReadyTimeout(opts.Cfg.TaskSchedule.Mailer.SendProcessor.ReadyTimeout),
			consume.WithStartReadDelay(opts.Cfg.TaskSchedule.Mailer.SendProcessor.StartReadDelay),
			consume.WithReadPeriod(opts.Cfg.TaskSchedule.Mailer.SendProcessor.ReadPeriod),
			consume.WithCancelReadTimeout(opts.Cfg.TaskSchedule.Mailer.SendProcessor.CancelReadTimeout),
			consume.WithHandlerTimeout(opts.Cfg.TaskSchedule.Mailer.SendProcessor.HandlerTimeout),
			consume.WithQueueSize(opts.Cfg.TaskSchedule.Mailer.SendProcessor.QueueSize),
			consume.WithWorkersCount(opts.Cfg.TaskSchedule.Mailer.SendProcessor.WorkersCount),
		),
		mrmailer.WithSendHandlerOpts(
			handle.WithClientEmail(mailProvider),
			handle.WithClientMessenger(telegramProvider),
		),
		mrmailer.WithTaskChangeFromToRetryOpts(
			task.WithCaption(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Caption),
			task.WithStartup(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Startup),
			task.WithPeriod(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeFromToRetry.Timeout),
		),
		mrmailer.WithTaskCleanMessagesOpts(
			task.WithCaption(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Caption),
			task.WithStartup(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Startup),
			task.WithPeriod(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Mailer.CleanQueue.Timeout),
		),
		mrmailer.WithChangeLimit(opts.Cfg.TaskSchedule.Mailer.ChangeQueueLimit),
		mrmailer.WithChangeRetryTimeout(opts.Cfg.TaskSchedule.Mailer.ChangeRetryTimeout),
		mrmailer.WithChangeRetryDelayed(opts.Cfg.TaskSchedule.Mailer.ChangeRetryDelayed),
		mrmailer.WithCleanLimit(opts.Cfg.TaskSchedule.Mailer.CleanQueueLimit),
	)

	return process, tasks, nil
}
