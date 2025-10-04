package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/pub/entity"
)

type (
	// MaterialType - comment struct.
	MaterialType struct {
		storage      pub.MaterialTypeStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewMaterialType - создаёт объект MaterialType.
func NewMaterialType(
	storage pub.MaterialTypeStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *MaterialType {
	return &MaterialType{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNameMaterialType),
	}
}

// GetList - comment method.
func (uc *MaterialType) GetList(ctx context.Context, lz mrcore.Localizer, params entity.MaterialTypeParams) ([]entity.MaterialType, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	if len(items) == 0 {
		return make([]entity.MaterialType, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
