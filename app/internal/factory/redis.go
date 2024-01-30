package factory

import (
	"context"
	"print-shop-back/config"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-webcore/mrlog"
)

func NewRedis(ctx context.Context, cfg config.Config) (*mrredis.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init redis connection")

	opts := mrredis.Options{
		Host:        cfg.Redis.Host,
		Port:        cfg.Redis.Port,
		Password:    cfg.Redis.Password,
		ConnTimeout: cfg.Redis.Timeout,
	}

	conn := mrredis.New()

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
