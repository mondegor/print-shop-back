package adm

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"

	"github.com/mondegor/print-shop-back/internal/controls/submitform/section/adm/entity"
)

type (
	// FormElementUseCase - comment interface.
	FormElementUseCase interface {
		GetItem(ctx context.Context, itemID uint64) (entity.FormElement, error)
		Create(ctx context.Context, item entity.FormElement) (itemID uint64, err error)
		Store(ctx context.Context, item entity.FormElement) error
		Remove(ctx context.Context, itemID uint64) error
		MoveAfterID(ctx context.Context, itemID, afterID uint64) error
	}

	// FormElementStorage - comment interface.
	FormElementStorage interface {
		NewCondition(formID uuid.UUID) mrstorage.SQLPartFunc
		Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormElement, error)
		FetchOne(ctx context.Context, rowID uint64) (entity.FormElement, error)
		FetchIDByParamName(ctx context.Context, formID uuid.UUID, paramName string) (rowID uint64, err error)
		IsExist(ctx context.Context, rowID uint64) error
		Insert(ctx context.Context, row entity.FormElement) (rowID uint64, err error)
		Update(ctx context.Context, row entity.FormElement) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
