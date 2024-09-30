package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/entity"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      pub.PrintFormatStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(storage pub.PrintFormatStorage, errorWrapper mrcore.UseCaseErrorWrapper) *PrintFormat {
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
