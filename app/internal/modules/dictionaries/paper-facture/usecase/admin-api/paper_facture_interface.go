package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	PaperFactureUseCase interface {
		GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.PaperFacture, error)
		Create(ctx context.Context, item entity.PaperFacture) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.PaperFacture) error
		ChangeStatus(ctx context.Context, item entity.PaperFacture) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	PaperFactureStorage interface {
		NewSelectParams(params entity.PaperFactureParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.PaperFacture, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.PaperFacture, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperFacture) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.PaperFacture) (int32, error)
		UpdateStatus(ctx context.Context, row entity.PaperFacture) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
