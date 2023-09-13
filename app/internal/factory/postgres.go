package factory

import (
    "context"
    "print-shop-back/config"
    "time"

    "github.com/mondegor/go-storage/mrpostgres"
    "github.com/mondegor/go-webcore/mrcore"
)

func NewPostgres(cfg *config.Config, logger mrcore.Logger) (*mrpostgres.ConnAdapter, error) {
    logger.Info("Create postgres connection")

    opt := mrpostgres.Options{
        Host: cfg.Storage.Host,
        Port: cfg.Storage.Port,
        Username: cfg.Storage.Username,
        Password: cfg.Storage.Password,
        Database: cfg.Storage.Database,
        MaxPoolSize: cfg.Storage.MaxPoolSize,
        ConnAttempts: 1,
        ConnTimeout: time.Duration(cfg.Storage.Timeout) * time.Second,
    }

    conn := mrpostgres.New()
    err := conn.Connect(opt)

    if err != nil {
        return nil, err
    }

    err = conn.Ping(context.TODO())

    return conn, err
}
