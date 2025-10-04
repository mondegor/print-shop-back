package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/module"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub"
	"github.com/mondegor/print-shop-back/internal/dictionaries/paperfacture/section/pub/entity"
)

type (
	// PaperFacture - comment struct.
	PaperFacture struct {
		storage      pub.PaperFactureStorage
		errorWrapper mrerr.UseCaseErrorWrapper
	}
)

// NewPaperFacture - создаёт объект PaperFacture.
func NewPaperFacture(
	storage pub.PaperFactureStorage,
	errorWrapper mrerr.UseCaseErrorWrapper,
) *PaperFacture {
	return &PaperFacture{
		storage:      storage,
		errorWrapper: mrerr.NewUseCaseErrorWrapper(errorWrapper, entity.ModelNamePaperFacture),
	}
}

// GetList - comment method.
func (uc *PaperFacture) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperFactureParams) ([]entity.PaperFacture, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err)
	}

	if len(items) == 0 {
		return make([]entity.PaperFacture, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
