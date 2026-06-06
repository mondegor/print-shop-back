package service

import (
	"github.com/mondegor/go-components/mrsettings"
	"github.com/mondegor/go-components/mrsettings/service/mustget"
	"github.com/mondegor/go-components/wire/mrsettings/cacheget"
	"github.com/mondegor/go-components/wire/mrsettings/dbset"
	"github.com/mondegor/go-sysmess/mrrun"
	"github.com/mondegor/go-sysmess/mrstorage/mrsql"
	"github.com/mondegor/go-sysmess/mrworker"
	"github.com/mondegor/go-sysmess/mrworker/job/task"

	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/app"
)

const (
	serviceSettingsTableName    = "printshop_global.settings"
	serviceSettingsPrimaryKey   = "setting_id"
	serviceSettingsTableNameLog = "printshop_global.settings_log"
)

// InitSettingsGetterAPI - создаёт получателя произвольных настроек из БД
// с использованием кэша и с периодическим его обновлением.
func InitSettingsGetterAPI(opts app.Options) (mrsettings.MustGetter, mrrun.Process) {
	log.Info(opts.Logger, "Create and init settings getter")

	getter, reloadScheduler := cacheget.InitServiceSettingsGetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		cacheget.WithCaptionPrefix("Settings/"),
		cacheget.WithTaskReloadSettingsOpts(
			task.WithCaptionPrefix("Settings/"),
			task.WithStartup(true),
			task.WithPeriodStrategy(
				mrworker.NewDoubleDelayedStartStrategy(
					opts.Cfg.TaskScheduleSettings.ReloadSettings.Period,
					opts.Cfg.TaskScheduleSettings.DefaultPeriodRatio,
				),
			),
			task.WithTimeout(opts.Cfg.TaskScheduleSettings.ReloadSettings.Timeout),
			task.WithSignalDo(
				opts.PostgresNotificationService.MustFind(opts.Cfg.TaskScheduleSettings.ReloadSettings.NotificationChannel),
			),
		),
	)

	return mustget.New(getter, opts.Logger), reloadScheduler
}

// InitSettingsSetterAPI - создаёт объект для сохранения произвольных настроек в БД.
func InitSettingsSetterAPI(opts app.Options) mrsettings.Setter {
	log.Info(opts.Logger, "Create and init settings setter")

	return dbset.InitServiceSettingsSetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		serviceSettingsTableNameLog,
	)
}
