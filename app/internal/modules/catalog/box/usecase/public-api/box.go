package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/box/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	Box struct {
		storage       BoxStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewBox(
	storage BoxStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *Box {
	return &Box{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	return items, nil
}
