package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/entity"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// Box - comment struct.
	// Box - comment struct.
	Box struct {
		storage      BoxStorage
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewBox - создаёт объект Box.
func NewBox(storage BoxStorage, errorWrapper mrcore.UsecaseErrorWrapper) *Box {
	return &Box{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *Box) GetList(ctx context.Context, params entity.BoxParams) ([]entity.Box, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameBox)
	}

	return items, nil
}
