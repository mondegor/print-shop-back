package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      PaperFactureStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(storage PaperFactureStorage, errorWrapper mrcore.UsecaseErrorWrapper) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *PaperFacture) GetList(ctx context.Context, params entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperFacture)
	}

	return items, nil
}
