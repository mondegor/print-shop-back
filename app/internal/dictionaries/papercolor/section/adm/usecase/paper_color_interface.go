package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PaperColorUseCase - comment interface.
	PaperColorUseCase interface {
		GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PaperColor, error)
		Create(ctx context.Context, item entity.PaperColor) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.PaperColor) error
		ChangeStatus(ctx context.Context, item entity.PaperColor) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		NewSelectParams(params entity.PaperColorParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.PaperColor, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperColor, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperColor) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.PaperColor) (int32, error)
		UpdateStatus(ctx context.Context, row entity.PaperColor) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
