package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// LaminateUseCase - comment interface.
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Laminate, error)
		Create(ctx context.Context, item entity.Laminate) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Laminate) error
		ChangeStatus(ctx context.Context, item entity.Laminate) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// LaminateStorage - comment interface.
	LaminateStorage interface {
		NewSelectParams(params entity.LaminateParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Laminate, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Laminate, error)
		FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Laminate) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Laminate) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Laminate) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
