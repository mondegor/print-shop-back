package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
)

type (
	// PrintFormatUseCase - comment interface.
	PrintFormatUseCase interface {
		GetList(ctx context.Context, params entity.PrintFormatParams) ([]entity.PrintFormat, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PrintFormat, error)
		Create(ctx context.Context, item entity.PrintFormat) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.PrintFormat) error
		ChangeStatus(ctx context.Context, item entity.PrintFormat) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		NewSelectParams(params entity.PrintFormatParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.PrintFormat, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PrintFormat, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PrintFormat) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.PrintFormat) (int32, error)
		UpdateStatus(ctx context.Context, row entity.PrintFormat) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
