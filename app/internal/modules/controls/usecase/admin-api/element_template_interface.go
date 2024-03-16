package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ElementTemplateUseCase interface {
		GetList(ctx context.Context, params entity.ElementTemplateParams) ([]entity.ElementTemplate, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.ElementTemplate, error)
		GetItemJson(ctx context.Context, itemID mrtype.KeyInt32, pretty bool) ([]byte, error)
		Create(ctx context.Context, item entity.ElementTemplate) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.ElementTemplate) error
		ChangeStatus(ctx context.Context, item entity.ElementTemplate) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
	}

	ElementTemplateStorage interface {
		NewFetchParams(params entity.ElementTemplateParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.ElementTemplate, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.ElementTemplate, error)
		FetchStatus(ctx context.Context, row entity.ElementTemplate) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.ElementTemplate) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.ElementTemplate) (int32, error)
		UpdateStatus(ctx context.Context, row entity.ElementTemplate) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
