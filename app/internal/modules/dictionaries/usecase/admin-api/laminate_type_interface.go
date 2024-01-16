package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	LaminateTypeService interface {
		GetList(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.LaminateType, error)
		Create(ctx context.Context, item *entity.LaminateType) error
		Store(ctx context.Context, item *entity.LaminateType) error
		ChangeStatus(ctx context.Context, item *entity.LaminateType) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	LaminateTypeStorage interface {
		NewFetchParams(params entity.LaminateTypeParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.LaminateType, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.LaminateType) error
		FetchStatus(ctx context.Context, row *entity.LaminateType) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.LaminateType) error
		Update(ctx context.Context, row *entity.LaminateType) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.LaminateType) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
