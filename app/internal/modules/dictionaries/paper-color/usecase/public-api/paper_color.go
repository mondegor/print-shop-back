package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-color/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	PaperColor struct {
		storage       PaperColorStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPaperColor(
	storage PaperColorStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperColor {
	return &PaperColor{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *PaperColor) GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	return items, nil
}
