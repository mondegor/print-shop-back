package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
)

type (
	// PaperUseCase - comment interface.
	PaperUseCase interface {
		GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Paper, error)
		Create(ctx context.Context, item entity.Paper) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Paper) error
		ChangeStatus(ctx context.Context, item entity.Paper) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// PaperStorage - comment interface.
	PaperStorage interface {
		NewSelectParams(params entity.PaperParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Paper, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Paper, error)
		FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Paper) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Paper) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Paper) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
