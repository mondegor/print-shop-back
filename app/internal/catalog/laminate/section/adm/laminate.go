package adm

import (
	"context"

	"print-shop-back/internal/adapter/workflow"
	"print-shop-back/internal/catalog/laminate/section/adm/entity"
)

type (
	// LaminateUseCase - comment interface.
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) (items []entity.Laminate, countItems int, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Laminate, error)
		Create(ctx context.Context, item entity.Laminate) (itemID uint64, err error)
		Save(ctx context.Context, item entity.Laminate) error
		ChangeStatus(ctx context.Context, item entity.Laminate) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// LaminateStorage - comment interface.
	LaminateStorage interface {
		FetchWithTotal(ctx context.Context, params entity.LaminateParams) (rows []entity.Laminate, countRows int, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Laminate, error)
		FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error)
		FetchStatus(ctx context.Context, rowID uint64) (workflow.ItemStatus, error)
		Insert(ctx context.Context, row entity.Laminate) (rowID uint64, err error)
		Update(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Laminate) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
