package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/entity/admin-api"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElementUseCase interface {
		GetList(ctx context.Context, params entity.FormElementParams) ([]entity.FormElement, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.FormElement, error)
		Create(ctx context.Context, item entity.FormElement) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.FormElement) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, itemID mrtype.KeyInt32, afterID mrtype.KeyInt32) error
	}

	FormElementStorage interface {
		GetMetaData(formID mrtype.KeyInt32) mrorderer.EntityMeta
		NewFetchParams(params entity.FormElementParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.FormElement, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.FormElement, error)
		FetchIdByName(ctx context.Context, formID mrtype.KeyInt32, paramName string) (mrtype.KeyInt32, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.FormElement) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.FormElement) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
