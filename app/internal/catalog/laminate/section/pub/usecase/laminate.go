package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/catalog/laminate/module"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub"
	"github.com/mondegor/print-shop-back/internal/catalog/laminate/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/measure"
)

type (
	// Laminate - comment struct.
	Laminate struct {
		storage      pub.LaminateStorage
		errorWrapper errors.Wrapper
	}
)

// NewLaminate - создаёт объект Laminate.
func NewLaminate(
	storage pub.LaminateStorage,
) *Laminate {
	return &Laminate{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetList - comment method.
func (uc *Laminate) GetList(ctx context.Context, lz mrcore.Localizer, params entity.LaminateParams) ([]entity.Laminate, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Laminate, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}

// GetTypeList - comment method.
func (uc *Laminate) GetTypeList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchTypeIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}

// GetThicknessList - comment method.
func (uc *Laminate) GetThicknessList(ctx context.Context) ([]measure.Meter, error) {
	items, err := uc.storage.FetchThicknesses(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}
