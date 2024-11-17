package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"
)

type (
	// BoxUseCase - comment interface.
	BoxUseCase interface {
		GetList(ctx context.Context, params entity.BoxParams) (items []entity.Box, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Box, error)
		Create(ctx context.Context, item entity.Box) (itemID uint64, err error)
		Store(ctx context.Context, item entity.Box) error
		ChangeStatus(ctx context.Context, item entity.Box) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// BoxStorage - comment interface.
	BoxStorage interface {
		FetchWithTotal(ctx context.Context, params entity.BoxParams) (rows []entity.Box, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Box, error)
		FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Box) (rowID uint64, err error)
		Update(ctx context.Context, row entity.Box) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Box) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
