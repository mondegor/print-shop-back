package service

import (
	"github.com/mondegor/go-components/mrauth/dto"
	"github.com/mondegor/go-components/wire/mrauth/scheduler"
	"github.com/mondegor/go-components/wire/mrauth/userstat/collector"
	"github.com/mondegor/go-core/mrprocess"
	"github.com/mondegor/go-core/mrprocess/collect"
	"github.com/mondegor/go-core/mrprocess/job/task"
	"github.com/mondegor/go-core/mrprocess/schedule"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
)

// TODO: дублирование название таблиц.
const (
	serviceAuthTokensTableName           = "printshop_auth.auth_tokens" //nolint:gosec
	serviceSecureOperationTableName      = "printshop_auth.secure_operations"
	serviceSecureOperationLogTableName   = "printshop_auth.secure_operations_log"
	serviceSessionsTableName             = "printshop_auth.sessions"
	serviceSessionsCleanupQueueTableName = "printshop_auth.sessions_cleanup_queue"
	serviceSessionsExcessQueueTableName  = "printshop_auth.sessions_excess_queue"
	// serviceUsersTableName              = "printshop_auth.users".
	serviceUsersActivityLogTableName  = "printshop_auth.users_activity_log"
	serviceUsersActivityStatTableName = "printshop_auth.users_activity_stat"
	// serviceUsersAuth2faTableName       = "printshop_auth.users_auth_2fa".
	// serviceUsersRealmsTableName        = "printshop_auth.users_realms".
)

// InitUserStatRequestCollectorService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitUserStatRequestCollectorService(opts app.Options) *collect.MessageCollector[dto.UserActivityLogMessage] {
	log.Info(opts.Logger, "Create and init user request collector service")

	return collector.NewService(
		opts.PostgresConnManager,
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		serviceUsersActivityLogTableName,
		serviceUsersActivityStatTableName,
		serviceSessionsTableName,
		collector.WithMessageCollectorOpts(
			collect.WithCaptionPrefix[dto.UserActivityLogMessage]("UserStat/"),
			collect.WithReadyTimeout[dto.UserActivityLogMessage](opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.ReadyTimeout),
			collect.WithFlushPeriodStrategy[dto.UserActivityLogMessage](
				mrprocess.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.FlushPeriod,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			collect.WithHandlerTimeout[dto.UserActivityLogMessage](opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.HandlerTimeout),
			collect.WithBatchSize[dto.UserActivityLogMessage](int(opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.BatchSize)),
			collect.WithWorkersCount[dto.UserActivityLogMessage](int(opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.WorkersCount)),
		),
	)
}

// InitAuthSchedulerService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitAuthSchedulerService(opts app.Options) *schedule.TaskScheduler {
	log.Info(opts.Logger, "Create and init auth scheduler service")

	return scheduler.NewService(
		opts.PostgresConnManager,
		opts.EventEmitter,
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		serviceAuthTokensTableName,
		serviceSecureOperationTableName,
		serviceSecureOperationLogTableName,
		serviceUsersActivityLogTableName,
		serviceSessionsTableName,
		serviceSessionsCleanupQueueTableName,
		serviceSessionsExcessQueueTableName,
		scheduler.WithCaptionPrefix("Auth/"),
		scheduler.WithCleanLimit(int(opts.Cfg.TaskScheduleAuth.CleanRecordsLimit)),
		scheduler.WithLogLifeTime(opts.Cfg.TaskScheduleAuth.LogsLifeTime),
		scheduler.WithTaskCleanRecordsOpts(
			task.WithCaptionPrefix("Auth/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrprocess.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleAuth.CleanRecords.Period,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskScheduleAuth.CleanRecords.Timeout),
		),
	)
}
