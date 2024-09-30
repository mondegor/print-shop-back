package adm

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// PaperFactureUseCase - comment interface.
	PaperFactureUseCase interface {
		GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PaperFacture, error)
		Create(ctx context.Context, item entity.PaperFacture) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.PaperFacture) error
		ChangeStatus(ctx context.Context, item entity.PaperFacture) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		NewSelectParams(params entity.PaperFactureParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.PaperFacture, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperFacture, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperFacture) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.PaperFacture) (int32, error)
		UpdateStatus(ctx context.Context, row entity.PaperFacture) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
