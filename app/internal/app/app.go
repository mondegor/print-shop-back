// Package app configures and runs application
package app

import (
    "context"
    "net/http"
    "print-shop-back/config"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/controller/http_v1"
    "print-shop-back/internal/factory"
    "print-shop-back/internal/infrastructure/repository"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "print-shop-back/pkg/mrhttp"
    "print-shop-back/pkg/mrlib"
    "time"

    sq "github.com/Masterminds/squirrel"
)

// Run creates objects via constructors.
func Run(cfg *config.Config, logger mrapp.Logger, translator mrapp.Translator) {
    appHelper := mrlib.NewHelper(logger)

    postgresClient, err := factory.NewPostgres(cfg, logger)
    appHelper.ExitOnError(err)
    defer appHelper.Close(postgresClient)

    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    requestValidator := mrlib.NewValidator()
    appHelper.ExitOnError(requestValidator.Register("article", dto.ValidateArticle))
    appHelper.ExitOnError(requestValidator.Register("variable", dto.ValidateVariable))

    errorHelper := mrerr.NewHelper()

    itemOrdererStorage := repository.NewItemOrderer(postgresClient, queryBuilder)
    itemOrdererComponent := usecase.NewItemOrdererComponent(itemOrdererStorage)

    catalogBoxStorage := repository.NewCatalogBox(postgresClient, queryBuilder)
    catalogBoxService := usecase.NewCatalogBox(catalogBoxStorage, errorHelper)
    catalogBoxHttp := http_v1.NewCatalogBox(catalogBoxService)

    catalogLaminateTypeStorage := repository.NewCatalogLaminateType(postgresClient, queryBuilder)
    catalogLaminateTypeService := usecase.NewCatalogLaminateType(catalogLaminateTypeStorage, errorHelper)
    catalogLaminateTypeHttp := http_v1.NewCatalogLaminateType(catalogLaminateTypeService)

    catalogLaminateStorage := repository.NewCatalogLaminate(postgresClient, queryBuilder)
    catalogLaminateService := usecase.NewCatalogLaminate(catalogLaminateStorage, catalogLaminateTypeStorage, errorHelper)
    catalogLaminateHttp := http_v1.NewCatalogLaminate(catalogLaminateService)

    catalogPaperColorStorage := repository.NewCatalogPaperColor(postgresClient, queryBuilder)
    catalogPaperColorService := usecase.NewCatalogPaperColor(catalogPaperColorStorage, errorHelper)
    catalogPaperColorHttp := http_v1.NewCatalogPaperColor(catalogPaperColorService)

    catalogPaperFactureStorage := repository.NewCatalogPaperFacture(postgresClient, queryBuilder)
    catalogPaperFactureService := usecase.NewCatalogPaperFacture(catalogPaperFactureStorage, errorHelper)
    catalogPaperFactureHttp := http_v1.NewCatalogPaperFacture(catalogPaperFactureService)

    catalogPaperStorage := repository.NewCatalogPaper(postgresClient, queryBuilder)
    catalogPaperService := usecase.NewCatalogPaper(catalogPaperStorage, catalogPaperColorStorage, catalogPaperFactureStorage, errorHelper)
    catalogPaperHttp := http_v1.NewCatalogPaper(catalogPaperService)

    catalogPrintFormatStorage := repository.NewCatalogPrintFormat(postgresClient, queryBuilder)
    catalogPrintFormatService := usecase.NewCatalogPrintFormat(catalogPrintFormatStorage, errorHelper)
    catalogPrintFormatHttp := http_v1.NewCatalogPrintFormat(catalogPrintFormatService)

    formFieldTemplateStorage := repository.NewFormFieldTemplate(postgresClient, queryBuilder)
    formFieldTemplateService := usecase.NewFormFieldTemplate(formFieldTemplateStorage, errorHelper)
    formFieldTemplateHttp := http_v1.NewFormFieldTemplate(formFieldTemplateService)

    formFieldItemStorage := repository.NewFormFieldItem(postgresClient, queryBuilder)

    formDataStorage := repository.NewFormData(postgresClient, queryBuilder)
    formDataService := usecase.NewFormData(formDataStorage, errorHelper)
    uiFormDataService := usecase.NewUIFormData(formDataStorage, formFieldItemStorage, errorHelper)
    formDataHttp := http_v1.NewFormData(formDataService, uiFormDataService)

    formFieldItemService := usecase.NewFormFieldItem(itemOrdererComponent, formFieldItemStorage, formFieldTemplateStorage, errorHelper)
    formFieldItemHttp := http_v1.NewFormFieldItem(formFieldItemService, formDataService, formFieldTemplateService)

    logger.Info("Create router")

    corsOptions := mrhttp.CorsOptions{
        AllowedOrigins: cfg.Cors.AllowedOrigins,
        AllowedMethods: cfg.Cors.AllowedMethods,
        AllowedHeaders: cfg.Cors.AllowedHeaders,
        ExposedHeaders: cfg.Cors.ExposedHeaders,
        AllowCredentials: cfg.Cors.AllowCredentials,
        Debug: cfg.Debug,
    }

    router := mrhttp.NewRouter(logger, requestValidator)
    router.RegisterMiddleware(
        mrhttp.NewCors(corsOptions),
        router.MiddlewareFirst(),
        router.MiddlewareAcceptLanguage(translator),
        router.MiddlewarePlatform(),
        router.MiddlewareAuthenticateUser(),
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

    server := mrhttp.NewServer(logger, mrhttp.ServerOptions{
        Handler: router,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 5 * time.Second,
        ShutdownTimeout: 30 * time.Second,
    })

    logger.Info("Start application")

    err = server.Start(mrhttp.ListenOptions{
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
        logger.Error(err)
    }
}

func MainPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"STATUS\": \"OK\"}"))
}
