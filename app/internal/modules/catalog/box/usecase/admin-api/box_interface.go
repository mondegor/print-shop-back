package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	BoxUseCase interface {
		GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Box, error)
		Create(ctx context.Context, item entity.Box) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Box) error
		ChangeStatus(ctx context.Context, item entity.Box) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	BoxStorage interface {
		NewSelectParams(params entity.BoxParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Box, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Box, error)
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Box) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Box) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Box) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
