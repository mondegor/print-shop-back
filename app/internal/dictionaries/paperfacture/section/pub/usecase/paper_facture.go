package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      pub.PaperFactureStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(storage pub.PaperFactureStorage, errorWrapper mrcore.UseCaseErrorWrapper) *PaperFacture {
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
