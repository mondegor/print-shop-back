package factory

import (
	"context"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/config"
)

// NewRedis - создаёт объект mrredis.ConnAdapter.
func NewRedis(ctx context.Context, cfg config.Config) (*mrredis.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init redis connection")

	opts := mrredis.Options{
		Host:         cfg.Redis.Host,
		Port:         cfg.Redis.Port,
		Password:     cfg.Redis.Password,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	}

	conn := mrredis.New()

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
