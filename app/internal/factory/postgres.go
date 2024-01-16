package factory

import (
	"context"
	"print-shop-back/config"
	"time"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewPostgres(cfg *config.Config, logger mrcore.Logger) (*mrpostgres.ConnAdapter, error) {
	logger.Info("Create and init postgres connection")

	opt := mrpostgres.Options{
		Host:         cfg.Storage.Host,
		Port:         cfg.Storage.Port,
		Username:     cfg.Storage.Username,
		Password:     cfg.Storage.Password,
		Database:     cfg.Storage.Database,
		MaxPoolSize:  cfg.Storage.MaxPoolSize,
		ConnAttempts: 1,
		ConnTimeout:  cfg.Storage.Timeout * time.Second,
	}

	conn := mrpostgres.New()

	if err := conn.Connect(opt); err != nil {
		return nil, err
	}

	return conn, conn.Ping(context.Background())
}
