package factory

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-sysmess/mrlog"

	"github.com/mondegor/print-shop-back/internal/app"
)

// InitRedis - создаёт объект mrredis.ConnAdapter.
func InitRedis(opts app.Options) (*mrredis.ConnAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg := opts.Cfg.Redis

	mrlog.Info(opts.Logger, "Create and init redis connection", "host", cfg.Host, "port", cfg.Port)

	redisOpts := mrredis.Options{
		Host:         cfg.Host,
		Port:         cfg.Port,
		Password:     cfg.Password,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	conn := mrredis.New(opts.Tracer)

	if err := conn.Connect(ctx, redisOpts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
