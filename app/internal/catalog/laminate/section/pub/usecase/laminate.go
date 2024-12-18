package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// Laminate - comment struct.
	Laminate struct {
		storage      pub.LaminateStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewLaminate - создаёт объект Laminate.
func NewLaminate(storage pub.LaminateStorage, errorWrapper mrcore.UseCaseErrorWrapper) *Laminate {
	return &Laminate{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *Laminate) GetList(ctx context.Context, params entity.LaminateParams) ([]entity.Laminate, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, nil
}

// GetTypeList - comment method.
func (uc *Laminate) GetTypeList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchTypeIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, nil
}

// GetThicknessList - comment method.
func (uc *Laminate) GetThicknessList(ctx context.Context) ([]measure.Meter, error) {
	items, err := uc.storage.FetchThicknesses(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNameLaminate)
	}

	return items, nil
}
