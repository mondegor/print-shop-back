package factory

import (
	"fmt"
	"io"
	"net/http"

	authiniting "github.com/mondegor/go-components/wire/mrauth/initing"
	"github.com/mondegor/go-storage/mrlock/redislocker"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xio"
	"github.com/mondegor/go-sysmess/wire"
	"github.com/mondegor/go-webcore/mrcore/initing"
	"github.com/mondegor/go-webcore/mrrun"

	dictionariesapi "github.com/mondegor/print-shop-back/cmd/factory/api/dictionaries"
	"github.com/mondegor/print-shop-back/cmd/factory/service"
	"github.com/mondegor/print-shop-back/cmd/factory/service/rest"
	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
)

// InitApp - Настраивает конфигурацию, внешнее окружение приложения, после этого создаёт её модули и компоненты.
func InitApp(args []string, stdout io.Writer) (app.Options, error) {
	parsedArgs, err := ParseAppArgs(args)
	if err != nil {
		return app.Options{}, err
	}

	cfg, err := config.Create(
		config.Args{
			WorkDir:     parsedArgs.WorkDir,
			ConfigPath:  parsedArgs.ConfigPath,
			DotEnvPath:  parsedArgs.DotEnvPath,
			Environment: parsedArgs.Environment,
			LogLevel:    parsedArgs.LogLevel,
			Stdout:      stdout,
		},
	)
	if err != nil {
		return app.Options{}, fmt.Errorf("factory.InitApp(): %w", err)
	}

	logger, tracer, err := InitLoggerAndTracer(cfg)
	if err != nil {
		return app.Options{}, err
	}

	traceManager, err := InitTraceContextManager(cfg, logger)
	if err != nil {
		return app.Options{}, err
	}

	return InitAppEnvironment(
		app.Options{
			Cfg:             cfg,
			Logger:          logger,
			Tracer:          tracer,
			TraceManager:    traceManager,
			OpenedResources: xio.NewCloseManager(logger),
		},
	)
}

// InitAppEnvironment - Настраивает внешнее окружение приложения на основе переданной конфигурации,
// после этого создаёт её модули и компоненты.
// Имеется возможность заранее задать некоторые параметры и компонентов приложения (актуально для использования в тестах):
// К ним относится: opts.PostgresConnManager, opts.RedisAdapter, opts.FileProviderPool.
func InitAppEnvironment(opts app.Options) (app.Options, error) {
	logger := opts.Logger

	// show head info about started app
	mrlog.Info(logger, opts.Cfg.App.Name, "environment", opts.Cfg.App.Environment, "version", opts.Cfg.App.Version)

	if opts.Cfg.Debugging.Debug {
		mrlog.Info(logger, "DEBUG MODE: ON")
	}

	mrlog.Info(logger, "LOG LEVEL: "+opts.Cfg.Log.Level)

	if opts.Cfg.App.WorkDir != "" {
		mrlog.Debug(logger, "WORK DIR: "+opts.Cfg.App.WorkDir)
	}

	mrlog.Debug(logger, "CONFIG PATH: "+opts.Cfg.ConfigPath)

	if opts.Cfg.App.DotEnvPath != "" {
		mrlog.Debug(logger, ".ENV PATH: "+opts.Cfg.App.DotEnvPath)
	}

	opts, err := createAppEnvironment(opts)
	if err != nil {
		return app.Options{}, err
	}

	// Shared APIs init section (!!! only after create environment)
	opts = createSharedAPI(opts)

	// Shared service's options (!!! only after create modules)
	if opts, err = createSharedServices(opts); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

// createAppEnvironment - создаёт, и настраивает внешнее окружение приложения.
func createAppEnvironment(opts app.Options) (enrichedOpts app.Options, err error) {
	if sentry, err := InitSentry(opts.Logger, opts.Cfg); err != nil {
		if !errors.Is(err, errSentryDisabled) {
			return app.Options{}, err
		}
	} else {
		opts.OpenedResources.Register(sentry)
		opts.Sentry = sentry
	}

	opts.InternalRouter = http.NewServeMux()

	if opts.Prometheus == nil {
		opts.Prometheus = InitPrometheus(opts)
	}

	// !!! only after init Sentry and Prometheus
	wire.InitErrors(
		wire.ErrorConfig{
			HasCaller:         opts.Cfg.Debugging.ErrorCaller.IsEnabled,
			CallerDepth:       opts.Cfg.Debugging.ErrorCaller.Depth,
			CallerShowFunc:    opts.Cfg.Debugging.ErrorCaller.ShowFunc,
			CallerUpperBounds: opts.Cfg.Debugging.ErrorCaller.UpperBounds,
		},
	)

	opts.EventEmitter = InitEventEmitter(opts)
	opts.ErrorHandler = wire.InitErrorHandler(opts.Logger)
	opts.AppHealth = mrrun.NewAppHealth()

	if opts.PostgresConnManager == nil {
		postgresAdapter, err := InitPostgres(opts)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources.Register(postgresAdapter)
		opts.PostgresConnManager = InitPostgresConnManager(postgresAdapter, opts.Logger)

		if opts.Cfg.Storage.MigrationsDir != "" {
			if err = ApplyPostgresMigrations(opts); err != nil {
				return app.Options{}, err
			}
		}
	}

	if opts.RedisAdapter == nil {
		redisAdapter, err := InitRedis(opts)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources.Register(redisAdapter)
		opts.RedisAdapter = redisAdapter
	}

	if opts.FileProviderPool == nil {
		opts.FileProviderPool, err = InitFileProviderPool(opts.Logger, opts.Tracer, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources.Register(opts.FileProviderPool)
	}

	redisCli, err := opts.RedisAdapter.Cli()
	if err != nil {
		return app.Options{}, err
	}

	opts.Locker = redislocker.NewLockerAdapter(
		redisCli,
		opts.Logger,
		opts.Tracer,
	)

	if opts.LocalePool, err = LocalePool(opts.Logger, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.RequestParsers, err = CreateRequestParsers(opts); err != nil {
		return app.Options{}, err
	}

	if opts.ResponseSenders, err = rest.CreateResponseSenders(opts.Logger); err != nil {
		return app.Options{}, err
	}

	if opts.PermsProvider, err = initing.InitPermsProvider(
		opts.Logger,
		opts.Cfg.AccessControl.RolesDirPath,
		opts.Cfg.AccessControl.Roles,
		opts.Cfg.AccessControl.Privileges,
		opts.Cfg.AccessControl.Permissions,
	); err != nil {
		return app.Options{}, err
	}

	opts.RealmUserProviders = authiniting.InitUserProviders(
		opts.Logger,
		opts.PostgresConnManager,
		authiniting.InitRealmKindRights(opts.Logger, opts.Cfg.Realms, opts.PermsProvider),
		opts.Cfg.Realms,
		opts.Cfg.Debugging.AuthorizedUser,
		opts.Cfg.AccessControl.JWTSecret,
	)

	if opts.ImageURLBuilder, err = InitImageURLBuilder(opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if err = RegisterSystemHandlers(opts); err != nil {
		return app.Options{}, err
	}

	if err = opts.Prometheus.Register(); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

func createSharedAPI(opts app.Options) app.Options {
	opts.PostgresNotificationService = mrpostgres.NewProcessWaitForNotification(
		opts.PostgresConnManager.ConnAdapter(),
		opts.Logger,
		[]string{
			opts.Cfg.TaskSchedule.Notifier.NoticeProcessor.NotificationChannel,
			opts.Cfg.TaskSchedule.Mailer.MessageProcessor.NotificationChannel,
			opts.Cfg.TaskSchedule.Settings.ReloadSettings.NotificationChannel,
		},
	)

	// create settings module
	{
		getter, reloadScheduler := service.InitSettingsGetterAPI(opts)
		opts.SettingsGetterAPI = getter
		opts.TaskSchedulerServices = append(opts.TaskSchedulerServices, reloadScheduler)

		opts.SettingsSetterAPI = service.InitSettingsSetterAPI(opts)
	}

	opts.MailerAPI = service.InitMailerAPI(opts)
	opts.NotifierAPI = service.InitNotifierAPI(opts)
	opts.DictionariesMaterialTypeAPI = dictionariesapi.InitMaterialTypeAvailabilityAPI(opts)
	opts.DictionariesPaperColorAPI = dictionariesapi.InitPaperColorAvailabilityAPI(opts)
	opts.DictionariesPaperFactureAPI = dictionariesapi.InitPaperFactureAvailabilityAPI(opts)
	opts.DictionariesPrintFormatAPI = dictionariesapi.InitPrintFormatAvailabilityAPI(opts)

	return opts
}

func createSharedServices(opts app.Options) (enrichedOpts app.Options, err error) {
	opts.UserStatRequestCollectorService = service.InitUserStatRequestCollectorService(opts)

	opts.MailProcessorService, err = service.InitMailerProcessorService(opts)
	if err != nil {
		return app.Options{}, fmt.Errorf("factory.InitMailerProcessorService(): %w", err)
	}

	opts.NoticeProcessorService = service.InitNotifierProcessorService(opts)

	opts.HttpServer, err = rest.InitRestServer(opts)
	if err != nil {
		return app.Options{}, fmt.Errorf("factory.InitRestServer(): %w", err)
	}

	opts.HttpInternalServer = InitInternalServer(opts)

	opts.TaskSchedulerServices = append(
		opts.TaskSchedulerServices,
		service.InitAuthSchedulerService(opts),
		service.InitMailerSchedulerService(opts),
		service.InitNotifierSchedulerService(opts),
	)

	return opts, nil
}
