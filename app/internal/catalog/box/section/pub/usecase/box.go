package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub"
	"github.com/mondegor/print-shop-back/internal/catalog/box/section/pub/entity"
)

type (
	// Box - comment struct.
	// Box - comment struct.
	Box struct {
		storage      pub.BoxStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewBox - создаёт объект Box.
func NewBox(storage pub.BoxStorage, errorWrapper mrcore.UseCaseErrorWrapper) *Box {
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
