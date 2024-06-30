package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/entity"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      PrintFormatStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(storage PrintFormatStorage, errorWrapper mrcore.UsecaseErrorWrapper) *PrintFormat {
	return &PrintFormat{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *PrintFormat) GetList(ctx context.Context, params entity.PrintFormatParams) ([]entity.PrintFormat, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePrintFormat)
	}

	return items, nil
}
