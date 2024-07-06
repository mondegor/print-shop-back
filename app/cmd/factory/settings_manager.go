package factory

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/app"

	"github.com/mondegor/go-components/factory/mrsettings"
	"github.com/mondegor/go-components/mrsettings/component/lightgetter"
	"github.com/mondegor/go-components/mrsettings/component/setter"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrworker/mrschedule"
)

const (
	settingsManagerTableName  = "printshop_global.settings"
	settingsManagerPrimaryKey = "setting_id"
)

// NewSettingsGetterAndTask - создаёт объекты: lightgetter.Component и mrschedule.TaskShell.
func NewSettingsGetterAndTask(ctx context.Context, opts app.Options) (*lightgetter.Component, *mrschedule.TaskShell) {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings getter")

	getter := mrsettings.NewComponentCacheGetter(
		opts.PostgresConnManager,
		mrsql.NewEntityMeta(settingsManagerTableName, settingsManagerPrimaryKey, nil),
		opts.UsecaseErrorWrapper,
		mrsettings.ComponentCacheGetterOptions{},
	)

	task := mrschedule.NewTaskShell(
		opts.Cfg.TaskSchedule.SettingsReloader.Caption,
		opts.Cfg.TaskSchedule.SettingsReloader.Startup,
		opts.Cfg.TaskSchedule.SettingsReloader.Period,
		opts.Cfg.TaskSchedule.SettingsReloader.Timeout,
		func(ctx context.Context) error {
			count, err := getter.Reload(ctx)
			if err != nil {
				return err
			}

			if count > 0 {
				mrlog.Ctx(ctx).Info().Msgf("Settings are reloaded: %d", count)
			}

			return nil
		},
	)

	return lightgetter.New(getter), task
}

// NewSettingsSetter - создаёт объект setter.Component.
func NewSettingsSetter(ctx context.Context, opts app.Options) *setter.Component {
	mrlog.Ctx(ctx).Info().Msg("Create and init settings setter")

	return mrsettings.NewComponentSetter(
		opts.PostgresConnManager,
		mrsql.NewEntityMeta(settingsManagerTableName, settingsManagerPrimaryKey, nil),
		opts.EventEmitter,
		opts.UsecaseErrorWrapper,
		mrsettings.ComponentSetterOptions{},
	)
}
