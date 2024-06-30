package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// BoxUseCase - comment interface.
	BoxUseCase interface {
		GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Box, error)
		Create(ctx context.Context, item entity.Box) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Box) error
		ChangeStatus(ctx context.Context, item entity.Box) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// BoxStorage - comment interface.
	BoxStorage interface {
		NewSelectParams(params entity.BoxParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Box, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Box, error)
		FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Box) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Box) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Box) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
