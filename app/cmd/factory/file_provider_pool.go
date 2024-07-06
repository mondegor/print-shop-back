package factory

import (
	"context"

	"github.com/mondegor/print-shop-back/config"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlog"
)

// NewFileProviderPool - создаёт объект mrstorage.FileProviderPool.
func NewFileProviderPool(ctx context.Context, cfg config.Config) (*mrstorage.FileProviderPool, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init file provider pool")

	pool := mrstorage.NewFileProviderPool()

	fs := NewFileSystem(ctx, cfg)

	if err := RegisterFileImageStorage(ctx, cfg, pool, fs); err != nil {
		return nil, err
	}

	return pool, nil
}
