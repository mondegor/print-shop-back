package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	LaminateType struct {
		storage       LaminateTypeStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewLaminateType(
	storage LaminateTypeStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *LaminateType {
	return &LaminateType{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *LaminateType) GetList(ctx context.Context, params entity.LaminateTypeParams) ([]entity.LaminateType, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminateType)
	}

	return items, nil
}
