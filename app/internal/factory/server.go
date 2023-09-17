package factory

import (
    "print-shop-back/config"
    "time"

    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrserver"
)

func NewHttpServer(cfg *config.Config, logger mrcore.Logger, router mrcore.HttpRouter) (*mrserver.ServerAdapter, error) {
    logger.Info("Create and start http server")

    server := mrserver.NewServer(logger, mrserver.ServerOptions{
        Handler: router,
        ReadTimeout: time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
        ShutdownTimeout: time.Duration(cfg.Server.ShutdownTimeout) * time.Second,
    })

    err := server.Start(mrserver.ListenOptions{
        AppPath: cfg.AppPath,
        Type: cfg.Listen.Type,
        SockName: cfg.Listen.SockName,
        BindIP: cfg.Listen.BindIP,
        Port: cfg.Listen.Port,
    })

    if err != nil {
        return nil, err
    }

    return server, nil
}
