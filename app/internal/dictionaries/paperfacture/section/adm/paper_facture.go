package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/adm/entity"
)

type (
	// PaperFactureUseCase - comment interface.
	PaperFactureUseCase interface {
		GetList(ctx context.Context, params entity.PaperFactureParams) (items []entity.PaperFacture, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.PaperFacture, error)
		Create(ctx context.Context, item entity.PaperFacture) (itemID uint64, err error)
		Store(ctx context.Context, item entity.PaperFacture) error
		ChangeStatus(ctx context.Context, item entity.PaperFacture) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// PaperFactureStorage - comment interface.
	PaperFactureStorage interface {
		FetchWithTotal(ctx context.Context, params entity.PaperFactureParams) (rows []entity.PaperFacture, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.PaperFacture, error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperFacture) (rowID uint64, err error)
		Update(ctx context.Context, row entity.PaperFacture) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.PaperFacture) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
