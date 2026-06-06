package service

import (
	"github.com/mondegor/go-components/mrauth/dto"
	"github.com/mondegor/go-components/wire/mrauth/scheduler"
	"github.com/mondegor/go-components/wire/mrauth/userstat/collector"
	"github.com/mondegor/go-sysmess/mrstorage/mrsql"
	"github.com/mondegor/go-sysmess/mrworker"
	"github.com/mondegor/go-sysmess/mrworker/job/task"
	"github.com/mondegor/go-sysmess/mrworker/process/collect"
	"github.com/mondegor/go-sysmess/mrworker/process/schedule"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
)

const (
	serviceAuthTokensTableName  = "printshop_auth.auth_tokens" //nolint:gosec
	serviceAuthTokensPrimaryKey = "refresh_token"

	serviceOperationTableName    = "printshop_auth.secure_operations"
	serviceOperationPrimaryKey   = "operation_token"
	serviceOperationLogTableName = "printshop_auth.secure_operations_log"

	serviceUserTableName     = "printshop_auth.users"
	serviceUserPrimaryKey    = "user_id"
	serviceUserStatTableName = "printshop_auth.users_activity_stat"
	serviceUserLogTableName  = "printshop_auth.users_activity_log"
)

// InitUserStatRequestCollectorService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitUserStatRequestCollectorService(opts app.Options) *collect.MessageCollector[dto.UserActivityLogMessage] {
	log.Info(opts.Logger, "Create and init user request collector service")

	return collector.NewService(
		opts.PostgresConnManager,
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceUserStatTableName,
			PrimaryKey: serviceUserPrimaryKey,
		},
		serviceUserLogTableName,
		collector.WithMessageCollectorOpts(
			collect.WithCaptionPrefix[dto.UserActivityLogMessage]("UserStat/"),
			collect.WithReadyTimeout[dto.UserActivityLogMessage](opts.Cfg.TaskScheduleAuth.UserStatRequestCollector.ReadyTimeout),
			collect.WithFlushPeriodStrategy[dto.UserActivityLogMessage](
				mrworker.NewDoubleDelayedStartStrategy(
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
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		mrsql.DBTableInfo{
			Name:       serviceAuthTokensTableName,
			PrimaryKey: serviceAuthTokensPrimaryKey,
		},
		mrsql.DBTableInfo{
			Name:       serviceOperationTableName,
			PrimaryKey: serviceOperationPrimaryKey,
		},
		serviceOperationLogTableName,
		serviceUserLogTableName,
		scheduler.WithCaptionPrefix("Auth/"),
		scheduler.WithCleanLimit(int(opts.Cfg.TaskScheduleAuth.CleanRecordsLimit)),
		scheduler.WithLogLifeTime(opts.Cfg.TaskScheduleAuth.LogsLifeTime),
		scheduler.WithTaskCleanRecordsOpts(
			task.WithCaptionPrefix("Auth/"),
			task.WithStartup(false),
			task.WithPeriodStrategy(
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleAuth.CleanRecords.Period,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskScheduleAuth.CleanRecords.Timeout),
		),
	)
}
