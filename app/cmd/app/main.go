package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "print-shop-back/config"
    "print-shop-back/internal/controller/http_v1"
    "print-shop-back/internal/controller/view"
    "print-shop-back/internal/factory"
    "print-shop-back/internal/infrastructure/repository"
    "print-shop-back/internal/usecase"
    "time"

    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-sysmess/mrlang"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrserver"
    "github.com/mondegor/go-webcore/mrtool"
    "github.com/mondegor/go-webcore/mrview"

    sq "github.com/Masterminds/squirrel"
)

const appName = "print-shop"
const appVersion = "v0.6.0"

var configPath string
var logLevel string

func init() {
   flag.StringVar(&configPath, "config-path", "./config/config.yaml", "Path to application config file")
   flag.StringVar(&logLevel, "log-level", "", "Logging level")
}

func main() {
    flag.Parse()

    cfg, err := config.New(configPath)

    if err != nil {
        log.Fatal(err)
    }

    if logLevel == "" {
        logLevel = cfg.Log.Level
    }

    logger, err := mrcore.NewLogger(appName, logLevel)

    if err != nil {
        log.Fatal(err)
    }

    logger.Info("APP VERSION: %s", appVersion)

    if cfg.Debug {
      logger.Info("DEBUG MODE: ON")
    }

    logger.Info("LOG LEVEL: %s", cfg.Log.Level)
    logger.Info("APP PATH: %s", cfg.AppPath)
    logger.Info("CONFIG PATH: %s", configPath)

    appHelper := mrtool.NewAppHelper(logger)

    responseTranslator, err := mrlang.NewTranslator(
        mrlang.TranslatorOptions{
            DirPath: cfg.Translation.DirPath,
            FileType: cfg.Translation.FileType,
            LangCodes: cfg.Translation.LangCodes,
        },
    )
    appHelper.ExitOnError(err)

    postgresClient, err := factory.NewPostgres(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(postgresClient)

    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    requestValidator := mrview.NewValidator()
    appHelper.ExitOnError(requestValidator.Register("article", view.ValidateArticle))
    appHelper.ExitOnError(requestValidator.Register("variable", view.ValidateVariable))

    serviceHelper := mrtool.NewServiceHelper()

    itemOrdererStorage := mrcom_orderer.NewRepository(postgresClient, queryBuilder)
    itemOrdererComponent := mrcom_orderer.NewComponent(itemOrdererStorage, logger)

    catalogBoxStorage := repository.NewCatalogBox(postgresClient, queryBuilder)
    catalogBoxService := usecase.NewCatalogBox(catalogBoxStorage, logger, serviceHelper)
    catalogBoxHttp := http_v1.NewCatalogBox(catalogBoxService)

    catalogLaminateTypeStorage := repository.NewCatalogLaminateType(postgresClient, queryBuilder)
    catalogLaminateTypeService := usecase.NewCatalogLaminateType(catalogLaminateTypeStorage, logger, serviceHelper)
    catalogLaminateTypeHttp := http_v1.NewCatalogLaminateType(catalogLaminateTypeService)

    catalogLaminateStorage := repository.NewCatalogLaminate(postgresClient, queryBuilder)
    catalogLaminateService := usecase.NewCatalogLaminate(catalogLaminateStorage, catalogLaminateTypeStorage, logger, serviceHelper)
    catalogLaminateHttp := http_v1.NewCatalogLaminate(catalogLaminateService)

    catalogPaperColorStorage := repository.NewCatalogPaperColor(postgresClient, queryBuilder)
    catalogPaperColorService := usecase.NewCatalogPaperColor(catalogPaperColorStorage, logger, serviceHelper)
    catalogPaperColorHttp := http_v1.NewCatalogPaperColor(catalogPaperColorService)

    catalogPaperFactureStorage := repository.NewCatalogPaperFacture(postgresClient, queryBuilder)
    catalogPaperFactureService := usecase.NewCatalogPaperFacture(catalogPaperFactureStorage, logger, serviceHelper)
    catalogPaperFactureHttp := http_v1.NewCatalogPaperFacture(catalogPaperFactureService)

    catalogPaperStorage := repository.NewCatalogPaper(postgresClient, queryBuilder)
    catalogPaperService := usecase.NewCatalogPaper(catalogPaperStorage, catalogPaperColorStorage, catalogPaperFactureStorage, logger, serviceHelper)
    catalogPaperHttp := http_v1.NewCatalogPaper(catalogPaperService)

    catalogPrintFormatStorage := repository.NewCatalogPrintFormat(postgresClient, queryBuilder)
    catalogPrintFormatService := usecase.NewCatalogPrintFormat(catalogPrintFormatStorage, logger, serviceHelper)
    catalogPrintFormatHttp := http_v1.NewCatalogPrintFormat(catalogPrintFormatService)

    formFieldTemplateStorage := repository.NewFormFieldTemplate(postgresClient, queryBuilder)
    formFieldTemplateService := usecase.NewFormFieldTemplate(formFieldTemplateStorage, logger, serviceHelper)
    formFieldTemplateHttp := http_v1.NewFormFieldTemplate(formFieldTemplateService)

    formFieldItemStorage := repository.NewFormFieldItem(postgresClient, queryBuilder)

    formDataStorage := repository.NewFormData(postgresClient, queryBuilder)
    formDataService := usecase.NewFormData(formDataStorage, logger, serviceHelper)
    uiFormDataService := usecase.NewUIFormData(formDataStorage, formFieldItemStorage, serviceHelper)
    formDataHttp := http_v1.NewFormData(formDataService, uiFormDataService)

    formFieldItemService := usecase.NewFormFieldItem(itemOrdererComponent, formFieldItemStorage, formFieldTemplateStorage, logger, serviceHelper)
    formFieldItemHttp := http_v1.NewFormFieldItem(formFieldItemService, formDataService, formFieldTemplateService)

    logger.Info("Create router")

    corsOptions := mrserver.CorsOptions{
        AllowedOrigins: cfg.Cors.AllowedOrigins,
        AllowedMethods: cfg.Cors.AllowedMethods,
        AllowedHeaders: cfg.Cors.AllowedHeaders,
        ExposedHeaders: cfg.Cors.ExposedHeaders,
        AllowCredentials: cfg.Cors.AllowCredentials,
        Debug: cfg.Debug,
    }

    router := mrserver.NewRouter(logger, mrserver.HandlerAdapter(requestValidator))
    router.RegisterMiddleware(
        mrserver.NewCors(corsOptions),
        mrserver.MiddlewareFirst(logger),
        mrserver.MiddlewareUserIp(),
        mrserver.MiddlewareAcceptLanguage(responseTranslator),
        mrserver.MiddlewarePlatform(mrcore.PlatformWeb),
        mrserver.MiddlewareAuthenticateUser(),
    )

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
    router.HandlerFunc(http.MethodGet, "/", MainPage)

    logger.Info("Initialize application")

    server := mrserver.NewServer(logger, mrserver.ServerOptions{
        Handler: router,
        ReadTimeout: time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
        ShutdownTimeout: time.Duration(cfg.Server.ShutdownTimeout) * time.Second,
    })

    logger.Info("Start application")

    err = server.Start(mrserver.ListenOptions{
        AppPath: cfg.AppPath,
        Type: cfg.Listen.Type,
        SockName: cfg.Listen.SockName,
        BindIP: cfg.Listen.BindIP,
        Port: cfg.Listen.Port,
    })
    appHelper.ExitOnError(err)
    defer appHelper.Close(server)

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go appHelper.GracefulShutdown(cancel)

    logger.Info("Waiting for requests. To exit press CTRL+C")

    select {
    case <-ctx.Done():
        err = server.Close()
        logger.Info("Application stopped")
    case err = <-server.Notify():
        logger.Info("Application stopped with error")
    }

    if err != nil && err != http.ErrServerClosed {
        logger.Err(err)
    }
}

func MainPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"STATUS\": \"OK\"}"))
}
