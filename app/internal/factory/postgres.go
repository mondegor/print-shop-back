package factory

import (
	"context"
	"print-shop-back/config"
	"time"

	"github.com/mondegor/go-storage/mrpostgres"
	"github.com/mondegor/go-webcore/mrlog"
)

func NewPostgres(ctx context.Context, cfg config.Config) (*mrpostgres.ConnAdapter, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init postgres connection")

	opts := mrpostgres.Options{
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

	if err := conn.Connect(ctx, opts); err != nil {
		return nil, err
	}

	return conn, conn.Ping(ctx)
}
