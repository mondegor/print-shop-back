package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateTypeUseCase interface {
		GetList(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.LaminateType, error)
		Create(ctx context.Context, item entity.LaminateType) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.LaminateType) error
		ChangeStatus(ctx context.Context, item entity.LaminateType) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	LaminateTypeStorage interface {
		NewSelectParams(params entity.LaminateTypeParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.LaminateType, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.LaminateType, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.LaminateType) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.LaminateType) (int32, error)
		UpdateStatus(ctx context.Context, row entity.LaminateType) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
