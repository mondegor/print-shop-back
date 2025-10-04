package factory

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-sysmess/mrerr/errorwrapper"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrwire"
	"github.com/mondegor/go-webcore/mrrun"

	"github.com/mondegor/print-shop-back/cmd/factory/auth"
	"github.com/mondegor/print-shop-back/cmd/factory/calculations"
	"github.com/mondegor/print-shop-back/cmd/factory/catalog"
	"github.com/mondegor/print-shop-back/cmd/factory/controls"
	"github.com/mondegor/print-shop-back/cmd/factory/dictionaries"
	"github.com/mondegor/print-shop-back/cmd/factory/filestation"
	"github.com/mondegor/print-shop-back/cmd/factory/provideraccounts"
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

	logger, err := InitLogger(cfg)
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
			Tracer:          InitTracer(cfg, logger),
			TraceManager:    traceManager,
			OpenedResources: mrwire.NewCloseManager(logger),
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
	if opts, err = createAppAPI(opts); err != nil {
		return app.Options{}, err
	}

	// Shared module's options (!!! only after create APIs)
	if opts, err = createAppModulesOptions(opts); err != nil {
		return app.Options{}, err
	}

	// Shared service's options (!!! only after create modules)
	if opts, err = createAppServices(opts); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

// createAppEnvironment - создаёт, и настраивает внешнее окружение приложения.
func createAppEnvironment(opts app.Options) (enrichedOpts app.Options, err error) {
	if opts.Cfg.Sentry.DSN != "" {
		sentry, err := InitSentry(opts.Logger, opts.Cfg)
		if err != nil {
			return app.Options{}, err
		}

		opts.OpenedResources.Register(sentry)
		opts.Sentry = sentry
	}

	opts.InternalRouter = http.NewServeMux()

	if opts.Prometheus == nil {
		opts.Prometheus = InitPrometheus(opts)
	}

	// !!! only after init Sentry and Prometheus
	InitProtoAppErrors(opts)

	opts.EventEmitter = InitEventEmitter(opts)
	opts.ErrorHandler = mrwire.InitErrorHandler(opts.Logger)
	opts.StorageErrorWrapper = errorwrapper.NewInfraStorage()
	opts.UsecaseErrorWrapper = errorwrapper.NewUseCase()
	opts.FileUserErrorWrapper = errorwrapper.NewDownloadUserImage()
	opts.ImageUserErrorWrapper = errorwrapper.NewDownloadUserImage()
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

	opts.Locker = mrredislock.NewLockerAdapter(redisCli, opts.Logger, opts.Tracer)

	if opts.LocalePool, err = LocalePool(opts.Logger, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	if opts.RequestParsers, err = CreateRequestParsers(opts); err != nil {
		return app.Options{}, err
	}

	if opts.ResponseSenders, err = rest.CreateResponseSenders(opts.Logger); err != nil {
		return app.Options{}, err
	}

	if opts.PermsProvider, err = InitPermsProvider(opts.Logger, opts.Cfg); err != nil {
		return app.Options{}, err
	}

	opts.RealmKindRights = InitRealmKindRights(opts.Logger, opts.Cfg.Realms, opts.PermsProvider)

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

func createAppAPI(opts app.Options) (enrichedOpts app.Options, err error) {
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

	if opts.DictionariesMaterialTypeAPI, err = dictionaries.NewMaterialTypeAvailabilityAPI(opts); err != nil {
		return app.Options{}, err
	}

	if opts.DictionariesPaperColorAPI, err = dictionaries.NewPaperColorAvailabilityAPI(opts); err != nil {
		return app.Options{}, err
	}

	if opts.DictionariesPaperFactureAPI, err = dictionaries.NewPaperFactureAvailabilityAPI(opts); err != nil {
		return app.Options{}, err
	}

	if opts.DictionariesPrintFormatAPI, err = dictionaries.NewPrintFormatAvailabilityAPI(opts); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

func createAppModulesOptions(opts app.Options) (enrichedOpts app.Options, err error) {
	if opts.AuthModule, err = auth.NewAuthModuleOptions(opts); err != nil {
		return app.Options{}, err
	}

	if opts.CalculationsAlgoModule, err = calculations.NewAlgoModuleOptions(opts); err != nil {
		return app.Options{}, err
	}

	if opts.CalculationsQueryHistoryModule, err = calculations.NewQueryHistoryModuleOptions(opts); err != nil {
		return app.Options{}, err
	}

	opts.CatalogBoxModule = catalog.NewBoxModuleOptions(opts)

	opts.CatalogLaminateModule = catalog.NewLaminateModuleOptions(opts)

	opts.CatalogPaperModule = catalog.NewPaperModuleOptions(opts)

	opts.ControlsElementTemplateModule = controls.NewElementTemplateModuleOptions(opts)

	opts.ControlsSubmitFormModule = controls.NewSubmitFormModuleOptions(opts)

	opts.DictionariesMaterialTypeModule = dictionaries.NewMaterialTypeModuleOptions(opts)

	opts.DictionariesPaperColorModule = dictionaries.NewPaperColorModuleOptions(opts)

	opts.DictionariesPaperFactureModule = dictionaries.NewPaperFactureModuleOptions(opts)

	opts.DictionariesPrintFormatModule = dictionaries.NewPrintFormatModuleOptions(opts)

	if opts.FileStationModule, err = filestation.NewModuleOptions(opts); err != nil {
		return app.Options{}, err
	}

	if opts.ProviderAccountsModule, err = provideraccounts.NewModuleOptions(opts); err != nil {
		return app.Options{}, err
	}

	return opts, nil
}

func createAppServices(opts app.Options) (enrichedOpts app.Options, err error) {
	opts.UserStatRequestCollectorService = service.InitUserStatRequestCollectorService(opts)

	opts.MailProcessorService, err = service.InitMailerProcessorService(opts)
	if err != nil {
		return app.Options{}, fmt.Errorf("factory.InitMailerService(): %w", err)
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
