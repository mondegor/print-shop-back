package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/catalog/paper/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	Paper struct {
		storage       PaperStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPaper(
	storage PaperStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *Paper {
	return &Paper{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *Paper) GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}
