package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// MaterialTypeUseCase - comment interface.
	MaterialTypeUseCase interface {
		GetList(ctx context.Context, params entity.MaterialTypeParams) ([]entity.MaterialType, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.MaterialType, error)
		Create(ctx context.Context, item entity.MaterialType) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.MaterialType) error
		ChangeStatus(ctx context.Context, item entity.MaterialType) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// MaterialTypeStorage - comment interface.
	MaterialTypeStorage interface {
		NewSelectParams(params entity.MaterialTypeParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.MaterialType, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.MaterialType, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.MaterialType) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.MaterialType) (int32, error)
		UpdateStatus(ctx context.Context, row entity.MaterialType) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
