package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
)

type (
	// PaperColorUseCase - comment interface.
	PaperColorUseCase interface {
		GetList(ctx context.Context, params entity.PaperColorParams) (items []entity.PaperColor, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.PaperColor, error)
		Create(ctx context.Context, item entity.PaperColor) (itemID uint64, err error)
		Store(ctx context.Context, item entity.PaperColor) error
		ChangeStatus(ctx context.Context, item entity.PaperColor) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// PaperColorStorage - comment interface.
	PaperColorStorage interface {
		FetchWithTotal(ctx context.Context, params entity.PaperColorParams) (rows []entity.PaperColor, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.PaperColor, error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.PaperColor) (rowID uint64, err error)
		Update(ctx context.Context, row entity.PaperColor) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.PaperColor) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
