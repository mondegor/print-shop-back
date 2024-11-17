package factory

import (
	"context"

	"github.com/mondegor/go-components/factory/mrsettings"
	"github.com/mondegor/go-components/factory/mrsettings/caching"
	"github.com/mondegor/go-components/mrsettings/component/lightget"
	"github.com/mondegor/go-components/mrsettings/component/set"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker"
	"github.com/mondegor/go-webcore/mrworker/job/task"

	"github.com/mondegor/print-shop-back/internal/app"
)

const (
	serviceSettingsTableName  = "printshop_global.settings"
	serviceSettingsPrimaryKey = "setting_id"
)

// NewSettingsGetterAPI - создаёт получателя произвольных настроек из БД
// с использованием кэша и с периодическим его обновлением.
func NewSettingsGetterAPI(ctx context.Context, opts app.Options) (*lightget.SettingsGetter, mrworker.Task) {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings getter")

	getter, reloadTask := caching.NewComponentGetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		caching.WithTaskReloadSettingsOpts(
			task.WithCaption(opts.Cfg.TaskSchedule.ReloadSettings.Caption),
			task.WithStartup(opts.Cfg.TaskSchedule.ReloadSettings.Startup),
			task.WithPeriod(opts.Cfg.TaskSchedule.ReloadSettings.Period),
			task.WithTimeout(opts.Cfg.TaskSchedule.ReloadSettings.Timeout),
		),
	)

	return lightget.New(getter), reloadTask
}

// NewSettingsSetterAPI - создаёт объект для сохранения произвольных настроек в БД.
func NewSettingsSetterAPI(ctx context.Context, opts app.Options) *set.SettingsSetter {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings setter")

	return mrsettings.NewComponentSetter(
		opts.PostgresConnManager,
		mrsql.DBTableInfo{
			Name:       serviceSettingsTableName,
			PrimaryKey: serviceSettingsPrimaryKey,
		},
		opts.EventEmitter,
	)
}
