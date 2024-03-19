package factory

import (
	"context"
	"log"
	"print-shop-back/config"
	"print-shop-back/internal"
	factory_catalog "print-shop-back/internal/factory/modules/catalog"
	factory_controls "print-shop-back/internal/factory/modules/controls"
	factory_dictionaries "print-shop-back/internal/factory/modules/dictionaries"
	factory_filestation "print-shop-back/internal/factory/modules/file-station"
	factory_provideraccounts "print-shop-back/internal/factory/modules/provider-accounts"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

func CreateAppEnvironment(configPath, logLevel string) (context.Context, app.Options) {
	cfg, err := config.Create(configPath)

	if err != nil {
		log.Fatal(err)
	}

	// create and init debugging
	mrdebug.SetDebugFlag(cfg.Debugging.Debug)

	mrerr.SetCallerOptions(
		mrerr.CallerDeep(cfg.Debugging.ErrorCaller.Deep),
		mrerr.CallerUseShortPath(cfg.Debugging.ErrorCaller.UseShortPath),
	)

	// create and init logger
	if logLevel != "" {
		cfg.Log.Level = logLevel
	}

	logger, err := NewLogger(cfg)

	if err != nil {
		log.Fatal(err)
	}

	mrlog.SetDefault(logger)

	// show head info about started app
	logger.Info().Msgf("%s, version: %s", cfg.AppName, cfg.AppVersion)

	if cfg.AppInfo != "" {
		logger.Info().Msg(cfg.AppInfo)
	}

	if mrdebug.IsDebug() {
		logger.Info().Msg("DEBUG MODE: ON")
	}

	logger.Info().Msgf("LOG LEVEL: %s", logger.Level())
	logger.Debug().Msgf("CONFIG PATH: %s", cfg.ConfigPath)
	logger.Debug().Msgf("APP PATH: %s", cfg.AppPath)

	ctx := mrlog.WithContext(context.Background(), logger)
	opts := app.Options{
		Cfg:          cfg,
		EventEmitter: logger,
	}

	return ctx, opts
}

func InitAppEnvironment(ctx context.Context, opts app.Options) (app.Options, error) {
	var err error

	// init shared options
	opts.UsecaseHelper = mrcore.NewUsecaseHelper()

	if opts.PostgresAdapter, err = NewPostgres(ctx, opts.Cfg); err != nil {
		return opts, err
	} else {
		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(opts.PostgresAdapter))
	}

	if opts.RedisAdapter, err = NewRedis(ctx, opts.Cfg); err != nil {
		return opts, err
	} else {
		opts.OpenedResources = append(opts.OpenedResources, mrlib.CloseFunc(opts.RedisAdapter))
	}

	if opts.FileProviderPool, err = NewFileProviderPool(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	opts.Locker = mrredislock.NewLockerAdapter(opts.RedisAdapter.Cli())

	if opts.Translator, err = NewTranslator(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.RequestParsers, err = CreateRequestParsers(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.ResponseSender, err = NewResponseSender(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	if opts.AccessControl, err = NewAccessControl(ctx, opts.Cfg); err != nil {
		return opts, err
	}

	opts.ImageURLBuilder = NewBuilderImagesURL(opts.Cfg)

	// Shared APIs init section (!!! only after init opts)
	if opts.DictionariesLaminateTypeAPI, err = factory_dictionaries.NewLaminateTypeAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperColorAPI, err = factory_dictionaries.NewPaperColorAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperFactureAPI, err = factory_dictionaries.NewPaperFactureAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPrintFormatAPI, err = factory_dictionaries.NewPrintFormatAPI(ctx, opts); err != nil {
		return opts, err
	}

	opts.OrdererAPI = NewOrdererAPI(ctx, opts)

	// Shared module's options (!!! only after init APIs)
	if opts.CatalogBoxModule, err = factory_catalog.NewBoxModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogLaminateModule, err = factory_catalog.NewLaminateModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.CatalogPaperModule, err = factory_catalog.NewPaperModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.ControlsElementTemplateModule, err = factory_controls.NewElementTemplateModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.ControlsSubmitFormModule, err = factory_controls.NewSubmitFormModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesLaminateTypeModule, err = factory_dictionaries.NewLaminateTypeModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperColorModule, err = factory_dictionaries.NewPaperColorModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperFactureModule, err = factory_dictionaries.NewPaperFactureModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPrintFormatModule, err = factory_dictionaries.NewPrintFormatModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.FileStationModule, err = factory_filestation.NewModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.ProviderAccountsModule, err = factory_provideraccounts.NewModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	return opts, nil
}
