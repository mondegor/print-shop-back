package adm

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/elementtemplate/section/adm/entity"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// ElementTemplateUseCase - comment interface.
	ElementTemplateUseCase interface {
		GetList(ctx context.Context, params entity.ElementTemplateParams) ([]entity.ElementTemplate, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.ElementTemplate, error)
		GetItemJson(ctx context.Context, itemID mrtype.KeyInt32, pretty bool) ([]byte, error)
		Create(ctx context.Context, item entity.ElementTemplate) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.ElementTemplate) error
		ChangeStatus(ctx context.Context, item entity.ElementTemplate) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	// ElementTemplateStorage - comment interface.
	ElementTemplateStorage interface {
		NewSelectParams(params entity.ElementTemplateParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.ElementTemplate, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplate, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.ElementTemplate) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.ElementTemplate) (int32, error)
		UpdateStatus(ctx context.Context, row entity.ElementTemplate) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
