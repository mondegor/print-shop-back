package factory

import (
	"context"
	"log"
	"print-shop-back/config"
	"print-shop-back/internal/modules"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrdebug"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
)

func CreateAppEnvironment(configPath, logLevel string) (context.Context, modules.Options) {
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
	opts := modules.Options{
		Cfg:          cfg,
		EventEmitter: logger,
	}

	return ctx, opts
}

func InitAppEnvironment(ctx context.Context, opts modules.Options) (modules.Options, error) {
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

	// Shared APIs init section (!!! only after init opts)
	if opts.DictionariesLaminateTypeAPI, err = NewDictionariesLaminateTypeAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperColorAPI, err = NewDictionariesPaperColorAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPaperFactureAPI, err = NewDictionariesPaperFactureAPI(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesPrintFormatAPI, err = NewDictionariesPrintFormatAPI(ctx, opts); err != nil {
		return opts, err
	}

	opts.OrdererAPI = NewOrdererAPI(ctx, opts)

	// Shared module's options (!!! only after init APIs)
	if opts.CatalogModule, err = NewCatalogModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.ControlsModule, err = NewControlsModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.DictionariesModule, err = NewDictionariesModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.FileStationModule, err = NewFileStationModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	if opts.ProviderAccountsModule, err = NewProviderAccountsModuleOptions(ctx, opts); err != nil {
		return opts, err
	}

	return opts, nil
}
