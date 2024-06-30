package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FormElementUseCase - comment interface.
	FormElementUseCase interface {
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.FormElement, error)
		Create(ctx context.Context, item entity.FormElement) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.FormElement) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, itemID, afterID mrtype.KeyInt32) error
	}

	// FormElementStorage - comment interface.
	FormElementStorage interface {
		NewOrderMeta(formID uuid.UUID) mrstorage.MetaGetter
		Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormElement, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.FormElement, error)
		FetchIDByParamName(ctx context.Context, formID uuid.UUID, paramName string) (mrtype.KeyInt32, error)
		IsExist(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.FormElement) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.FormElement) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
