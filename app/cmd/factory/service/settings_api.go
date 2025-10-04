package service

import (
	"github.com/mondegor/go-components/factory/mrsettings"
	"github.com/mondegor/go-components/factory/mrsettings/caching"
	"github.com/mondegor/go-components/mrsettings/usecase/liteget"
	"github.com/mondegor/go-components/mrsettings/usecase/set"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-webcore/mrworker/job/task"
	"github.com/mondegor/go-webcore/mrworker/process/schedule"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceSettingsTableName    = "printshop_global.settings"
	serviceSettingsPrimaryKey   = "setting_id"
	serviceSettingsTableNameLog = "printshop_global.settings_log"
)

// InitSettingsGetterAPI - создаёт получателя произвольных настроек из БД
// с использованием кэша и с периодическим его обновлением.
func InitSettingsGetterAPI(opts app.Options) (*liteget.SettingsGetter, *schedule.TaskScheduler) {
	mrlog.Info(opts.Logger, "Create and init settings getter")

	getter, reloadScheduler := caching.NewComponentGetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		opts.ErrorHandler,
		opts.Logger,
		opts.TraceManager,
		caching.WithCaptionPrefix("Settings/"),
		caching.WithTaskReloadSettingsOpts(
			task.WithCaptionPrefix("Settings/"),
			task.WithStartup(true),
			task.WithPeriod(opts.Cfg.TaskSchedule.Settings.ReloadSettings.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.Settings.ReloadSettings.Timeout),
			task.WithSignalDo(
				opts.PostgresNotificationService.ReceiverChannels.MustFind(opts.Cfg.TaskSchedule.Settings.ReloadSettings.NotificationChannel),
			),
		),
	)

	return liteget.New(getter, opts.Logger), reloadScheduler
}

// InitSettingsSetterAPI - создаёт объект для сохранения произвольных настроек в БД.
func InitSettingsSetterAPI(opts app.Options) *set.SettingsSetter {
	mrlog.Info(opts.Logger, "Create and init settings setter")

	return mrsettings.NewComponentSetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		serviceSettingsTableNameLog,
		opts.EventEmitter,
	)
}
