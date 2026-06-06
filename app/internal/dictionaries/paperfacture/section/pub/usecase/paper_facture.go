package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrcore"

	"print-shop-back/internal/dictionaries/paperfacture/module"
	"print-shop-back/internal/dictionaries/paperfacture/section/pub"
	"print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      pub.PaperFactureStorage
		errorWrapper errors.Wrapper
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	storage pub.PaperFactureStorage,
) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetList - comment method.
func (uc *PaperFacture) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.PaperFacture, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
