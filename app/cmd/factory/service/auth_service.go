package service

import (
	"github.com/mondegor/go-components/mrauth/dto"
	"github.com/mondegor/go-components/wire/mrauth/scheduler"
	"github.com/mondegor/go-components/wire/mrauth/userstat/collector"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/collect"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/internal/app"
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
	mrlog.Info(opts.Logger, "Create and init user request collector service")

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
			collect.WithReadyTimeout[dto.UserActivityLogMessage](opts.Cfg.TaskSchedule.Auth.UserStat.RequestCollector.ReadyTimeout),
			collect.WithFlushPeriod[dto.UserActivityLogMessage](opts.Cfg.TaskSchedule.Auth.UserStat.RequestCollector.FlushPeriod),
			collect.WithHandlerTimeout[dto.UserActivityLogMessage](opts.Cfg.TaskSchedule.Auth.UserStat.RequestCollector.HandlerTimeout),
			collect.WithBatchSize[dto.UserActivityLogMessage](int(opts.Cfg.TaskSchedule.Auth.UserStat.RequestCollector.BatchSize)),
			collect.WithWorkersCount[dto.UserActivityLogMessage](int(opts.Cfg.TaskSchedule.Auth.UserStat.RequestCollector.WorkersCount)),
		),
	)
}

// InitAuthSchedulerService - создаёт сервис для обработки сообщений и связанных с ним задачи.
func InitAuthSchedulerService(opts app.Options) *schedule.TaskScheduler {
	mrlog.Info(opts.Logger, "Create and init auth scheduler service")

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
		scheduler.WithCleanLimit(int(opts.Cfg.TaskSchedule.Auth.CleanRecordsLimit)),
		scheduler.WithLogLifeTime(opts.Cfg.TaskSchedule.Auth.LogsLifeTime),
		scheduler.WithTaskCleanRecordsOpts(
			task.WithCaptionPrefix("Auth/"),
			task.WithStartup(false),
			task.WithPeriod(opts.Cfg.TaskSchedule.Auth.CleanRecords.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Auth.CleanRecords.Timeout),
		),
	)
}
