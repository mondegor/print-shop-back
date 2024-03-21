package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/laminate/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	Laminate struct {
		storage       LaminateStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewLaminate(
	storage LaminateStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *Laminate {
	return &Laminate{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *Laminate) GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, nil
}
