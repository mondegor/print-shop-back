package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/box/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	BoxService interface {
		GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Box, error)
		Create(ctx context.Context, item *entity.Box) error
		Store(ctx context.Context, item *entity.Box) error
		ChangeStatus(ctx context.Context, item *entity.Box) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	BoxStorage interface {
		NewFetchParams(params entity.BoxParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Box, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Box) error
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row *entity.Box) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.Box) error
		Update(ctx context.Context, row *entity.Box) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.Box) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
