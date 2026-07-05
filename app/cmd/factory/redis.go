package factory

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrredis"

	"print-shop-back/config"
	"print-shop-back/internal/adapter/log"
	"print-shop-back/internal/adapter/trace"
)

// InitRedis - создаёт объект mrredis.ConnAdapter.
func InitRedis(logger log.Logger, tracer trace.Tracer, cfg config.Config) (*mrredis.ConnAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Info(logger, "Create and init redis connection", "host", cfg.RedisHost, "port", cfg.RedisPort)

	redisOpts := mrredis.Options{
		Host:         cfg.RedisHost,
		Port:         cfg.RedisPort,
		Password:     cfg.RedisPassword,
		ReadTimeout:  cfg.RedisReadTimeout,
		WriteTimeout: cfg.RedisWriteTimeout,
	}

	conn := mrredis.New(tracer)

	if err := conn.Connect(ctx, redisOpts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
