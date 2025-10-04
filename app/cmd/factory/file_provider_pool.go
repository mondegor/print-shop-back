package factory

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtrace"

	"github.com/mondegor/print-shop-back/config"
)

// InitFileProviderPool - создаёт объект mrstorage.FileProviderPool.
func InitFileProviderPool(logger mrlog.Logger, tracer mrtrace.Tracer, cfg config.Config) (*mrstorage.FileProviderPool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	mrlog.Info(logger, "Create and init file provider pool")

	pool := mrstorage.NewFileProviderPool()

	fsAdapter := NewFileSystem(logger, cfg)

	if err := RegisterFileImageStorage(logger, tracer, cfg, pool, fsAdapter); err != nil {
		return nil, err
	}

	return pool, pool.Ping(ctx)
}
