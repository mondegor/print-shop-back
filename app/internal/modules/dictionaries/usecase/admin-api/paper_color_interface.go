package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperColorService interface {
		GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.PaperColor, error)
		Create(ctx context.Context, item *entity.PaperColor) error
		Store(ctx context.Context, item *entity.PaperColor) error
		ChangeStatus(ctx context.Context, item *entity.PaperColor) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	PaperColorStorage interface {
		NewFetchParams(params entity.PaperColorParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PaperColor, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.PaperColor) error
		FetchStatus(ctx context.Context, row *entity.PaperColor) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.PaperColor) error
		Update(ctx context.Context, row *entity.PaperColor) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.PaperColor) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
