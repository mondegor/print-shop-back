package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "print-shop-back/config"
    "print-shop-back/internal/controller/http_v1"
    "print-shop-back/internal/factory"
    "print-shop-back/internal/infrastructure/repository"
    "print-shop-back/internal/usecase"

    sq "github.com/Masterminds/squirrel"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-webcore/mrtool"
)

var (
    configPath string
    logLevel string
)

func init() {
    flag.StringVar(&configPath,"config-path", "./config/config.yaml", "Path to application config file")
    flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
    flag.Parse()

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

    appHelper := mrtool.NewAppHelper(logger)
    serviceHelper := mrtool.NewServiceHelper()

    postgresAdapter, err := factory.NewPostgres(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(postgresAdapter)

    // redisAdapter, err := factory.NewRedis(cfg, logger)
    // appHelper.ExitOnError(err)
    // defer appHelper.Close(redisAdapter)

    // lockerAdapter := mrredsync.NewLockerAdapter(redisAdapter.Cli())
    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    itemOrdererStorage := mrcom_orderer.NewRepository(postgresAdapter, queryBuilder)
    itemOrdererComponent := mrcom_orderer.NewComponent(itemOrdererStorage, logger)

    catalogBoxStorage := repository.NewCatalogBox(postgresAdapter, queryBuilder)
    catalogBoxService := usecase.NewCatalogBox(catalogBoxStorage, logger, serviceHelper)
    catalogBoxHttp := http_v1.NewCatalogBox(catalogBoxService)

    catalogLaminateTypeStorage := repository.NewCatalogLaminateType(postgresAdapter, queryBuilder)
    catalogLaminateTypeService := usecase.NewCatalogLaminateType(catalogLaminateTypeStorage, logger, serviceHelper)
    catalogLaminateTypeHttp := http_v1.NewCatalogLaminateType(catalogLaminateTypeService)

    catalogLaminateStorage := repository.NewCatalogLaminate(postgresAdapter, queryBuilder)
    catalogLaminateService := usecase.NewCatalogLaminate(catalogLaminateStorage, catalogLaminateTypeStorage, logger, serviceHelper)
    catalogLaminateHttp := http_v1.NewCatalogLaminate(catalogLaminateService)

    catalogPaperColorStorage := repository.NewCatalogPaperColor(postgresAdapter, queryBuilder)
    catalogPaperColorService := usecase.NewCatalogPaperColor(catalogPaperColorStorage, logger, serviceHelper)
    catalogPaperColorHttp := http_v1.NewCatalogPaperColor(catalogPaperColorService)

    catalogPaperFactureStorage := repository.NewCatalogPaperFacture(postgresAdapter, queryBuilder)
    catalogPaperFactureService := usecase.NewCatalogPaperFacture(catalogPaperFactureStorage, logger, serviceHelper)
    catalogPaperFactureHttp := http_v1.NewCatalogPaperFacture(catalogPaperFactureService)

    catalogPaperStorage := repository.NewCatalogPaper(postgresAdapter, queryBuilder)
    catalogPaperService := usecase.NewCatalogPaper(catalogPaperStorage, catalogPaperColorStorage, catalogPaperFactureStorage, logger, serviceHelper)
    catalogPaperHttp := http_v1.NewCatalogPaper(catalogPaperService)

    catalogPrintFormatStorage := repository.NewCatalogPrintFormat(postgresAdapter, queryBuilder)
    catalogPrintFormatService := usecase.NewCatalogPrintFormat(catalogPrintFormatStorage, logger, serviceHelper)
    catalogPrintFormatHttp := http_v1.NewCatalogPrintFormat(catalogPrintFormatService)

    formFieldTemplateStorage := repository.NewFormFieldTemplate(postgresAdapter, queryBuilder)
    formFieldTemplateService := usecase.NewFormFieldTemplate(formFieldTemplateStorage, logger, serviceHelper)
    formFieldTemplateHttp := http_v1.NewFormFieldTemplate(formFieldTemplateService)

    formFieldItemStorage := repository.NewFormFieldItem(postgresAdapter, queryBuilder)

    formDataStorage := repository.NewFormData(postgresAdapter, queryBuilder)
    formDataService := usecase.NewFormData(formDataStorage, logger, serviceHelper)
    uiFormDataService := usecase.NewUIFormData(formDataStorage, formFieldItemStorage, serviceHelper)
    formDataHttp := http_v1.NewFormData(formDataService, uiFormDataService)

    formFieldItemService := usecase.NewFormFieldItem(itemOrdererComponent, formFieldItemStorage, formFieldTemplateStorage, logger, serviceHelper)
    formFieldItemHttp := http_v1.NewFormFieldItem(formFieldItemService, formDataService, formFieldTemplateService)

    router, err := factory.NewHttpRouter(cfg, logger)
    appHelper.ExitOnError(err)

    router.Register(
        catalogBoxHttp,
        catalogLaminateHttp,
        catalogLaminateTypeHttp,
        catalogPaperHttp,
        catalogPaperColorHttp,
        catalogPaperFactureHttp,
        catalogPrintFormatHttp,
        formFieldTemplateHttp,
        formDataHttp,
        formFieldItemHttp,
    )

    serverAdapter, err := factory.NewHttpServer(cfg, logger, router)
    appHelper.ExitOnError(err)
    defer appHelper.Close(serverAdapter)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go appHelper.GracefulShutdown(cancel)

    logger.Info("Waiting for requests. To exit press CTRL+C")

    select {
    case <-ctx.Done():
        err = serverAdapter.Close()
        logger.Info("Application stopped")
    case err = <-serverAdapter.Notify():
        logger.Info("Application stopped with error")
    }

    if err != nil && err != http.ErrServerClosed {
        logger.Err(err)
    }
}
