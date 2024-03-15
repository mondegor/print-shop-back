package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateUseCase interface {
		GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Laminate, error)
		Create(ctx context.Context, item entity.Laminate) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Laminate) error
		ChangeStatus(ctx context.Context, item entity.Laminate) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	LaminateStorage interface {
		NewFetchParams(params entity.LaminateParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Laminate, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Laminate, error)
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row entity.Laminate) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.Laminate) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Laminate) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Laminate) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
