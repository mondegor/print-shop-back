package usecase

import (
	"context"

	"github.com/mondegor/go-core/errors"
	"github.com/mondegor/go-webcore/mrcore"

	"print-shop-back/internal/dictionaries/papercolor/module"
	"print-shop-back/internal/dictionaries/papercolor/section/pub"
	"print-shop-back/internal/dictionaries/papercolor/section/pub/entity"
)

type (
	// PaperColor - comment struct.
	PaperColor struct {
		storage      pub.PaperColorStorage
		errorWrapper errors.Wrapper
	}
)

// NewPaperColor - создаёт объект PaperColor.
func NewPaperColor(
	storage pub.PaperColorStorage,
) *PaperColor {
	return &PaperColor{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetList - comment method.
func (uc *PaperColor) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperColorParams) ([]entity.PaperColor, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.PaperColor, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
