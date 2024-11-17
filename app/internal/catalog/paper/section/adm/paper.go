package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/adm/entity"
)

type (
	// PaperUseCase - comment interface.
	PaperUseCase interface {
		GetList(ctx context.Context, params entity.PaperParams) (items []entity.Paper, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Paper, error)
		Create(ctx context.Context, item entity.Paper) (itemID uint64, err error)
		Store(ctx context.Context, item entity.Paper) error
		ChangeStatus(ctx context.Context, item entity.Paper) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// PaperStorage - comment interface.
	PaperStorage interface {
		FetchWithTotal(ctx context.Context, params entity.PaperParams) (rows []entity.Paper, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Paper, error)
		FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Paper) (id uint64, err error)
		Update(ctx context.Context, row entity.Paper) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Paper) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
