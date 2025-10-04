package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/pub/entity"
)

type (
	// PrintFormat - comment struct.
	PrintFormat struct {
		storage      pub.PrintFormatStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewPrintFormat - создаёт объект PrintFormat.
func NewPrintFormat(
	storage pub.PrintFormatStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *PrintFormat {
	return &PrintFormat{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNamePrintFormat),
	}
}

// GetList - comment method.
func (uc *PrintFormat) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PrintFormatParams) ([]entity.PrintFormat, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	if len(items) == 0 {
		return make([]entity.PrintFormat, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
