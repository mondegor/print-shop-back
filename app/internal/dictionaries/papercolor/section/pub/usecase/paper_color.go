package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      pub.PaperColorStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(storage pub.PaperColorStorage, errorWrapper mrcore.UseCaseErrorWrapper) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *PaperColor) GetList(ctx context.Context, params entity.PaperColorParams) ([]entity.PaperColor, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaperColor)
	}

	return items, nil
}
