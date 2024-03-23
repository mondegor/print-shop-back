package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/controls/submit-form/entity/admin-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FormElementUseCase interface {
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.FormElement, error)
		Create(ctx context.Context, item entity.FormElement) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.FormElement) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, itemID mrtype.KeyInt32, afterID mrtype.KeyInt32) error
	}

	FormElementStorage interface {
		NewOrderMeta(formID uuid.UUID) mrorderer.EntityMeta
		Fetch(ctx context.Context, formID uuid.UUID) ([]entity.FormElement, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.FormElement, error)
		FetchIdByParamName(ctx context.Context, formID uuid.UUID, paramName string) (mrtype.KeyInt32, error)
		IsExist(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.FormElement) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.FormElement) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
