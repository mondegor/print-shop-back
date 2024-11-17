package adm

import (
	"context"

	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"
)

type (
	// ElementTemplateUseCase - comment interface.
	ElementTemplateUseCase interface {
		GetList(ctx context.Context, params entity.ElementTemplateParams) (items []entity.ElementTemplate, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.ElementTemplate, error)
		GetItemJson(ctx context.Context, itemID uint64, pretty bool) ([]byte, error)
		Create(ctx context.Context, item entity.ElementTemplate) (itemID uint64, err error)
		Store(ctx context.Context, item entity.ElementTemplate) error
		ChangeStatus(ctx context.Context, item entity.ElementTemplate) error
		Remove(ctx context.Context, itemID uint64) error
	}

	// ElementTemplateStorage - comment interface.
	ElementTemplateStorage interface {
		FetchWithTotal(ctx context.Context, params entity.ElementTemplateParams) (rows []entity.ElementTemplate, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.ElementTemplate, error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.ElementTemplate) (rowID uint64, err error)
		Update(ctx context.Context, row entity.ElementTemplate) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.ElementTemplate) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
