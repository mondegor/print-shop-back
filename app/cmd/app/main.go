package main

import (
	"context"
	"flag"
	"log"
	"print-shop-back/config"
	"print-shop-back/internal/factory"
	"print-shop-back/internal/modules"
	factory_catalog_adm "print-shop-back/internal/modules/catalog/factory/admin-api"
	factory_controls_adm "print-shop-back/internal/modules/controls/factory/admin-api"
	factory_dictionaries_adm "print-shop-back/internal/modules/dictionaries/factory/admin-api"
	factory_filestation_pub "print-shop-back/internal/modules/file-station/factory/public-api"
	factory_provider_accounts_adm "print-shop-back/internal/modules/provider-accounts/factory/admin-api"
	factory_provider_accounts_pacc "print-shop-back/internal/modules/provider-accounts/factory/provider-account-api"
	factory_provider_account_pub "print-shop-back/internal/modules/provider-accounts/factory/public-api"

	"github.com/mondegor/go-storage/mrredislock"
	"github.com/mondegor/go-webcore/mrtool"
)

var (
	configPath string
	logLevel   string
)

func init() {
	flag.StringVar(&configPath, "config-path", "./config/config.yaml", "Path to application config file")
	flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
	flag.Parse()

	sharedOptions := &modules.Options{}
	cfg, err := config.New(configPath)

	if err != nil {
		log.Fatal(err)
	}

	if logLevel != "" {
		cfg.Log.Level = logLevel
	}

	logger, err := factory.NewLogger(cfg)

	if err != nil {
		log.Fatal(err)
	}

	sharedOptions.Cfg = cfg
	sharedOptions.Logger = logger
	sharedOptions.EventBox = logger

	appHelper := mrtool.NewAppHelper(logger)
	sharedOptions.ServiceHelper = mrtool.NewServiceHelper()

	sharedOptions.PostgresAdapter, err = factory.NewPostgres(cfg, logger)
	appHelper.ExitOnError(err)
	defer appHelper.Close(sharedOptions.PostgresAdapter)

	sharedOptions.RedisAdapter, err = factory.NewRedis(cfg, logger)
	appHelper.ExitOnError(err)
	defer appHelper.Close(sharedOptions.RedisAdapter)

	sharedOptions.FileProviderPool, err = factory.NewFileProviderPool(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.Locker = mrredislock.NewLockerAdapter(sharedOptions.RedisAdapter.Cli())
	sharedOptions.OrdererAPI = factory.NewOrdererAPI(cfg, sharedOptions.PostgresAdapter, logger, logger)

	sharedOptions.Translator, err = factory.NewTranslator(cfg, logger)
	appHelper.ExitOnError(err)

	sharedOptions.DictionariesLaminateTypeAPI, err = factory.NewDictionariesLaminateTypeAPI(sharedOptions)
	appHelper.ExitOnError(err)

	sharedOptions.DictionariesPaperColorAPI, err = factory.NewDictionariesPaperColorAPI(sharedOptions)
	appHelper.ExitOnError(err)

	sharedOptions.DictionariesPaperFactureAPI, err = factory.NewDictionariesPaperFactureAPI(sharedOptions)
	appHelper.ExitOnError(err)

	sharedOptions.DictionariesPrintFormatAPI, err = factory.NewDictionariesPrintFormatAPI(sharedOptions)
	appHelper.ExitOnError(err)

	// module's access
	modulesAccess, err := factory.NewModulesAccess(cfg, logger)
	appHelper.ExitOnError(err)

	// module's options
	catalogOptions, err := factory.NewCatalogOptions(sharedOptions)
	appHelper.ExitOnError(err)

	controlsOptions, err := factory.NewControlsOptions(sharedOptions)
	appHelper.ExitOnError(err)

	dictionariesAreaOptions, err := factory.NewDictionariesOptions(sharedOptions)
	appHelper.ExitOnError(err)

	fileStationOptions, err := factory.NewFileStationOptions(sharedOptions)
	appHelper.ExitOnError(err)

	providerAccountsOptions, err := factory.NewProviderAccountsOptions(sharedOptions)
	appHelper.ExitOnError(err)

	// http router
	router, err := factory.NewHttpRouter(cfg, logger, sharedOptions.Translator)
	appHelper.ExitOnError(err)

	// section: admin-api
	sectionAdminAPI := factory.NewClientSectionAdminAPI(cfg, logger, modulesAccess)

	appHelper.ExitOnError(
		factory.RegisterSystemHandlers(cfg, logger, router, sectionAdminAPI),
	)

	controllers, err := factory_catalog_adm.NewModule(catalogOptions, sectionAdminAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	controllers, err = factory_controls_adm.NewModule(controlsOptions, sectionAdminAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	controllers, err = factory_dictionaries_adm.NewModule(dictionariesAreaOptions, sectionAdminAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	controllers, err = factory_provider_accounts_adm.NewModule(providerAccountsOptions, sectionAdminAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	// section: provider-account-api
	sectionProviderAccountAPI := factory.NewClientSectionProviderAccountAPI(cfg, logger, modulesAccess)

	appHelper.ExitOnError(
		factory.RegisterSystemHandlers(cfg, logger, router, sectionProviderAccountAPI),
	)

	controllers, err = factory_provider_accounts_pacc.NewModule(providerAccountsOptions, sectionProviderAccountAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	// section: public-api
	sectionPublicAPI := factory.NewClientSectionPublicAPI(cfg, logger, modulesAccess)

	appHelper.ExitOnError(
		factory.RegisterSystemHandlers(cfg, logger, router, sectionPublicAPI),
	)

	controllers, err = factory_filestation_pub.NewModule(fileStationOptions, sectionPublicAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	controllers, err = factory_provider_account_pub.NewModule(providerAccountsOptions, sectionPublicAPI)
	appHelper.ExitOnError(err)
	router.Register(controllers...)

	// http server
	serverAdapter, err := factory.NewHttpServer(cfg, logger, router)
	appHelper.ExitOnError(err)
	defer appHelper.Close(serverAdapter)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go appHelper.GracefulShutdown(cancel)

	logger.Info("Waiting for requests. To exit press CTRL+C")

	select {
	case <-ctx.Done():
		logger.Info("Application stopped")
	case err = <-serverAdapter.Notify():
		logger.Info("Application stopped with error")
	}

	logger.Err(err)
}
