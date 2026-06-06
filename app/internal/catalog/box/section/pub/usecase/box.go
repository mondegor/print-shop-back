package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrcore"

	"print-shop-back/internal/catalog/box/module"
	"print-shop-back/internal/catalog/box/section/pub"
	"print-shop-back/internal/catalog/box/section/pub/entity"
)

type (
	// Box - comment struct.
	Box struct {
		storage      pub.BoxStorage
		errorWrapper errors.Wrapper
	}
)

// NewBox - создаёт объект Box.
func NewBox(
	storage pub.BoxStorage,
) *Box {
	return &Box{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetList - comment method.
func (uc *Box) GetList(ctx context.Context, lz mrcore.Localizer, params entity.BoxParams) ([]entity.Box, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Box, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}
