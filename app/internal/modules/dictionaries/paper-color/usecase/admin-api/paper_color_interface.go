package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColorUseCase interface {
		GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PaperColor, error)
		Create(ctx context.Context, item entity.PaperColor) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.PaperColor) error
		ChangeStatus(ctx context.Context, item entity.PaperColor) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	PaperColorStorage interface {
		NewSelectParams(params entity.PaperColorParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PaperColor, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperColor, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperColor) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.PaperColor) (int32, error)
		UpdateStatus(ctx context.Context, row entity.PaperColor) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
