package factory

import (
	"context"
	"print-shop-back/config"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-webcore/mrcore"
)

func NewRedis(cfg *config.Config, logger mrcore.Logger) (*mrredis.ConnAdapter, error) {
	logger.Info("Create and init redis connection")

	opt := mrredis.Options{
		Host:        cfg.Redis.Host,
		Port:        cfg.Redis.Port,
		Password:    cfg.Redis.Password,
		ConnTimeout: cfg.Redis.Timeout,
	}

	conn := mrredis.New()

	if err := conn.Connect(opt); err != nil {
		return nil, err
	}

	return conn, conn.Ping(context.Background())
}
