package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"
)

type (
	// LaminateUseCase - comment interface.
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) (items []entity.Laminate, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Laminate, error)
		Create(ctx context.Context, item entity.Laminate) (itemID uint64, err error)
		Store(ctx context.Context, item entity.Laminate) error
		ChangeStatus(ctx context.Context, item entity.Laminate) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// LaminateStorage - comment interface.
	LaminateStorage interface {
		FetchWithTotal(ctx context.Context, params entity.LaminateParams) (rows []entity.Laminate, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Laminate, error)
		FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Laminate) (rowID uint64, err error)
		Update(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
