package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
)

type (
	// MaterialTypeUseCase - comment interface.
	MaterialTypeUseCase interface {
		GetList(ctx context.Context, params entity.MaterialTypeParams) (items []entity.MaterialType, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.MaterialType, error)
		Create(ctx context.Context, item entity.MaterialType) (itemID uint64, err error)
		Store(ctx context.Context, item entity.MaterialType) error
		ChangeStatus(ctx context.Context, item entity.MaterialType) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		FetchWithTotal(ctx context.Context, params entity.MaterialTypeParams) (rows []entity.MaterialType, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.MaterialType, error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.MaterialType) (rowID uint64, err error)
		Update(ctx context.Context, row entity.MaterialType) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.MaterialType) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
