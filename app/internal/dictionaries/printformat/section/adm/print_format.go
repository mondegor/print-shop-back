package adm

import (
	"context"

	"github.com/mondegor/go-sysmess/mrstatus/itemstatus"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
)

type (
	// PrintFormatUseCase - comment interface.
	PrintFormatUseCase interface {
		GetList(ctx context.Context, params entity.PrintFormatParams) (items []entity.PrintFormat, countItems int, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.PrintFormat, error)
		Create(ctx context.Context, item entity.PrintFormat) (itemID uint64, err error)
		Save(ctx context.Context, item entity.PrintFormat) error
		ChangeStatus(ctx context.Context, item entity.PrintFormat) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// PrintFormatStorage - comment interface.
	PrintFormatStorage interface {
		FetchWithTotal(ctx context.Context, params entity.PrintFormatParams) (rows []entity.PrintFormat, countRows int, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.PrintFormat, error)
		FetchStatus(ctx context.Context, rowID uint64) (itemstatus.Enum, error)
		Insert(ctx context.Context, row entity.PrintFormat) (rowID uint64, err error)
		Update(ctx context.Context, row entity.PrintFormat) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.PrintFormat) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
