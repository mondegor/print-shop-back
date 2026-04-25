package factory

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	authcfg "github.com/mondegor/go-components/wire/mrauth/config"
	authiniting "github.com/mondegor/go-components/wire/mrauth/initing"
	"github.com/mondegor/go-storage/mrlock/redislocker"
	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/xio"
	"github.com/mondegor/go-sysmess/wire"
	"github.com/mondegor/go-webcore/mrrun"
	wirecore "github.com/mondegor/go-webcore/wire"

	dictionariesapi "github.com/mondegor/print-shop-back/cmd/factory/api/dictionaries"
	"github.com/mondegor/print-shop-back/cmd/factory/service"
	"github.com/mondegor/print-shop-back/cmd/factory/service/rest"
	"github.com/mondegor/print-shop-back/config"
	"github.com/mondegor/print-shop-back/internal/app"
)

// InitApp - Настраивает конфигурацию, внешнее окружение приложения, после этого создаёт её модули и компоненты.
func InitApp(args []string, stdout io.Writer) (app.Options, error) {
	parsedArgs, err := config.ParseCmdArgs(args)
	if err != nil {
		return app.Options{}, err
	}

	cfg, err := config.Create(parsedArgs, stdout)
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
			DebugFunc:       InitDebugInfo(cfg.DebugIsEnabled),
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
	mrlog.Info(logger, opts.Cfg.AppName, "environment", opts.Cfg.Environment, "version", opts.Cfg.AppVersion)

	if opts.Cfg.WorkDir != "" {
		mrlog.Debug(logger, "WORK DIR: "+opts.Cfg.WorkDir)
	}

	mrlog.Debug(logger, "CONFIG PATHs: "+strings.Join(opts.Cfg.ConfigPaths, ", "))

	if opts.Cfg.DotEnvPath != "" {
		mrlog.Debug(logger, ".ENV PATH: "+opts.Cfg.DotEnvPath)
	}

	if opts.Cfg.DebugIsEnabled {
		mrlog.Info(logger, "DEBUG MODE: ON")
	}

	mrlog.Info(logger, "LOG LEVEL: "+opts.Cfg.LogLevel)

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

	opts.MonitoringRouter = http.NewServeMux()

	if opts.Prometheus == nil {
		opts.Prometheus = InitPrometheus(opts)
	}

	// !!! only after init Sentry and Prometheus
	wire.InitErrors(
		wire.ErrorConfig{
			HasCaller:         opts.Cfg.StackTraceIsEnabled,
			CallerDepth:       opts.Cfg.StackTraceDepth,
			CallerShowFunc:    opts.Cfg.StackTraceShowFunc,
			CallerUpperBounds: opts.Cfg.StackTraceUpperBounds,
		},
	)

	opts.EventEmitter = InitEventEmitter(opts)
	opts.ErrorHandler = wire.InitErrorHandler(opts.Logger)
	opts.AppHealth = mrrun.NewAppHealth()

	if opts.PostgresConnManager == nil {
		postgresAdapter, err := InitPostgres(opts.Logger, opts.Tracer, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources.Register(postgresAdapter)
		opts.PostgresConnManager = InitPostgresConnManager(postgresAdapter, opts.Logger)

		if err = InitPrometheusStatPostgres(opts); err != nil {
			return app.Options{}, err
		}

		if err = ApplyPostgresMigrations(opts.Logger, opts.PostgresConnManager, opts.Cfg); err != nil {
			return app.Options{}, err
		}
	}

	if opts.RedisAdapter == nil {
		redisAdapter, err := InitRedis(opts.Logger, opts.Tracer, opts.Cfg)
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

	if opts.PermsProvider, err = wirecore.InitPermsProvider(
		opts.Logger,
		opts.Cfg.AccessControl.RolesDirPath,
		opts.Cfg.AccessControl.Roles,
		opts.Cfg.AccessControl.AllowedPrivileges,
		opts.Cfg.AccessControl.AllowedPermissions,
	); err != nil {
		return app.Options{}, err
	}

	rights, err := authiniting.InitRealmKindRights(opts.Logger, opts.Cfg.AccessControl.Realms, opts.PermsProvider)
	if err != nil {
		return app.Options{}, err
	}

	opts.RealmUserProviders = authiniting.InitUserProviders(
		opts.Logger,
		opts.PostgresConnManager,
		rights,
		opts.Cfg.AccessControl.Realms,
		authcfg.TestUser{
			ID:       opts.Cfg.TestUserID,
			Realm:    opts.Cfg.TestUserRealm,
			Kind:     opts.Cfg.TestUserKind,
			LangCode: opts.Cfg.TestUserLangCode,
		},
		opts.Cfg.AccessControl.JWTSecret,
	)

	if opts.ImageURLBuilder, err = InitImageURLBuilder(opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if err = RegisterSystemHandlers(opts); err != nil {
		return app.Options{}, err
	}

	if opts.Prometheus != nil {
		if err = opts.Prometheus.Register(); err != nil {
			return app.Options{}, err
		}
	}

	return opts, nil
}

func createSharedAPI(opts app.Options) app.Options {
	opts.PostgresNotificationService = mrpostgres.NewProcessWaitForNotification(
		opts.PostgresConnManager.ConnAdapter(),
		opts.Logger,
		[]string{
			opts.Cfg.TaskScheduleNotifier.NoticeProcessor.NotificationChannel,
			opts.Cfg.TaskScheduleMailer.MessageProcessor.NotificationChannel,
			opts.Cfg.TaskScheduleSettings.ReloadSettings.NotificationChannel,
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

	opts.HttpMonitoringServer = InitMonitoringServer(opts)

	opts.TaskSchedulerServices = append(
		opts.TaskSchedulerServices,
		service.InitAuthSchedulerService(opts),
		service.InitMailerSchedulerService(opts),
		service.InitNotifierSchedulerService(opts),
	)

	return opts, nil
}
