package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      pub.MaterialTypeStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(storage pub.MaterialTypeStorage, errorWrapper mrcore.UseCaseErrorWrapper) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *MaterialType) GetList(ctx context.Context, params entity.MaterialTypeParams) ([]entity.MaterialType, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameMaterialType)
	}

	return items, nil
}
