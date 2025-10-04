package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      pub.PaperColorStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage pub.PaperColorStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNamePaperColor),
	}
}

// GetList - comment method.
func (uc *PaperColor) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperColorParams) ([]entity.PaperColor, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	if len(items) == 0 {
		return make([]entity.PaperColor, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
