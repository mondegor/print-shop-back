package usecase

import (
	"context"
	entity "print-shop-back/internal/modules/dictionaries/paper-facture/entity/public-api"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	PaperFacture struct {
		storage       PaperFactureStorage
		usecaseHelper *mrcore.UsecaseHelper
	}
)

func NewPaperFacture(
	storage PaperFactureStorage,
	usecaseHelper *mrcore.UsecaseHelper,
) *PaperFacture {
	return &PaperFacture{
		storage:       storage,
		usecaseHelper: usecaseHelper,
	}
}

func (uc *PaperFacture) GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	items, err := uc.storage.Fetch(ctx, params)

	if err != nil {
		return nil, uc.usecaseHelper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	return items, nil
}
