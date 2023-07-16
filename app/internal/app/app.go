// Package app configures and runs application
package app

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/signal"
    "print-shop-back/config"
    "print-shop-back/internal/controller/dto"
    "print-shop-back/internal/controller/http_v1"
    "print-shop-back/internal/infrastructure/repository"
    "print-shop-back/internal/usecase"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrapp"
    "print-shop-back/pkg/mrerr"
    "print-shop-back/pkg/mrhttp"
    "print-shop-back/pkg/mrlib"
    "syscall"
    "time"

    sq "github.com/Masterminds/squirrel"
)

// Run creates objects via constructors.
func Run(cfg *config.Config, logger mrapp.Logger, translator mrapp.Translator) {
    //logger.Info("Create redis connection: %s:%s", cfg.Redis.Host, cfg.Redis.Port)
    //
    //redisOptions := mrredis.Options{
    //    Host: cfg.Redis.Host,
    //    Port: cfg.Redis.Port,
    //    Password: cfg.Redis.Password,
    //    ConnTimeout: time.Duration(cfg.Redis.Timeout),
    //}
    //
    //redisClient := mrredis.New(logger)
    //
    //if err := redisClient.Connect(redisOptions); err != nil {
    //    logger.Fatal("%s", err.Error())
    //}

    logger.Info("Create postgres connection: %s:%s, DB=%s", cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.Database)

    postgresOptions := mrpostgres.Options{
        Host: cfg.Storage.Host,
        Port: cfg.Storage.Port,
        Username: cfg.Storage.Username,
        Password: cfg.Storage.Password,
        Database: cfg.Storage.Database,
        MaxPoolSize: 1,
        ConnAttempts: 1,
        ConnTimeout: time.Duration(cfg.Storage.Timeout),
    }

    postgresClient := mrpostgres.New()
    queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

    if err := postgresClient.Connect(context.TODO(), postgresOptions); err != nil {
        logger.Fatal("%s", err.Error())
    }

    requestValidator := mrlib.NewValidator()
    appExitIfError(logger, requestValidator.Register("article", dto.ValidateArticle))
    appExitIfError(logger, requestValidator.Register("variable", dto.ValidateVariable))

    errorHelper := mrerr.NewHelper()

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

    formFieldItemOrdererStorage := repository.NewFormFieldItemOrderer(postgresClient, queryBuilder)
    formFieldItemService := usecase.NewFormFieldItem(formFieldItemStorage, formFieldTemplateStorage, errorHelper)
    formFieldItemOrdererService := usecase.NewFormFieldItemOrderer(formFieldItemOrdererStorage, errorHelper)
    formFieldItemHttp := http_v1.NewFormFieldItem(formFieldItemService, formFieldItemOrdererService, formDataService, formFieldTemplateService)

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

    appStart(cfg, logger, router)
}

func appStart(cfg *config.Config, logger mrapp.Logger, router mrapp.Router) {
    logger.Info("Initialize application")

    server := mrhttp.NewServer(logger, mrhttp.ServerOptions{
        Handler: router,
        ReadTimeout: 5 * time.Second,
        WriteTimeout: 5 * time.Second,
        ShutdownTimeout: 3 * time.Second,
    })

    logger.Info("Start application")

    server.Start(mrhttp.ListenOptions{
        AppPath: cfg.AppPath,
        Type: cfg.Listen.Type,
        SockName: cfg.Listen.SockName,
        BindIP: cfg.Listen.BindIP,
        Port: cfg.Listen.Port,
    })

    signalAppChan := make(chan os.Signal, 1)
    signal.Notify(
        signalAppChan,
        syscall.SIGABRT,
        syscall.SIGQUIT,
        syscall.SIGHUP,
        os.Interrupt,
        syscall.SIGTERM,
    )

    select {
        case signalApp := <-signalAppChan:
            logger.Info("Application shutdown, signal: " + signalApp.String())

        case err := <-server.Notify():
            logger.Error(fmt.Errorf("http server shutdown: %w", err))
    }

    closeItems := []io.Closer{server}

    for _, closer := range closeItems {
        if err := closer.Close(); err != nil {
            logger.Error(fmt.Errorf("failed to close %v: %w", closer, err))
        }
    }
}

func appExitIfError(logger mrapp.Logger, err error) {
    if err != nil {
        logger.Fatal(err)
    }
}

func MainPage(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("{\"STATUS\": \"OK\"}"))
}
