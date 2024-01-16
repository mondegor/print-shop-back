package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperService interface {
		GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Paper, error)
		Create(ctx context.Context, item *entity.Paper) error
		Store(ctx context.Context, item *entity.Paper) error
		ChangeStatus(ctx context.Context, item *entity.Paper) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	PaperStorage interface {
		NewFetchParams(params entity.PaperParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Paper, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Paper) error
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row *entity.Paper) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.Paper) error
		Update(ctx context.Context, row *entity.Paper) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.Paper) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
