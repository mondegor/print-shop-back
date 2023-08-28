package factory

import (
    "context"
    "print-shop-back/config"
    "print-shop-back/pkg/client/mrpostgres"
    "print-shop-back/pkg/mrapp"
    "time"
)

func NewPostgres(cfg *config.Config, logger mrapp.Logger) (*mrpostgres.Connection, error) {
    logger.Info("Create postgres connection")

    opt := mrpostgres.Options{
        Host: cfg.Storage.Host,
        Port: cfg.Storage.Port,
        Username: cfg.Storage.Username,
        Password: cfg.Storage.Password,
        Database: cfg.Storage.Database,
        MaxPoolSize: 1,
        ConnAttempts: 1,
        ConnTimeout: time.Duration(cfg.Storage.Timeout),
    }

    conn := mrpostgres.New()
    err := conn.Connect(context.TODO(), opt)

    return conn, err
}
