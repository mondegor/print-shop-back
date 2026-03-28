package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/module"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/measure"
)

type (
	// Paper - comment struct.
	Paper struct {
		storage      pub.PaperStorage
		errorWrapper errors.Wrapper
	}
)

// NewPaper - создаёт объект Paper.
func NewPaper(
	storage pub.PaperStorage,
) *Paper {
	return &Paper{
		storage:      storage,
		errorWrapper: errors.NewServiceRecordNotFoundWrapper(),
	}
}

// GetList - comment method.
func (uc *Paper) GetList(ctx context.Context, lz mrcore.Localizer, params entity.PaperParams) ([]entity.Paper, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	if len(items) == 0 {
		return make([]entity.Paper, 0), nil
	}

	for i := range items {
		items[i].Caption = lz.TranslateInDomain(module.LocaleDomain, items[i].Caption)
	}

	return items, nil
}

// GetTypeList - comment method.
func (uc *Paper) GetTypeList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchTypeIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}

// GetColorList - comment method.
func (uc *Paper) GetColorList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchColorIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}

// GetDensityList - comment method.
func (uc *Paper) GetDensityList(ctx context.Context) ([]measure.KilogramPerMeter2, error) {
	items, err := uc.storage.FetchDensities(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}

// GetFactureList - comment method.
func (uc *Paper) GetFactureList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchFactureIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.Wrap(err)
	}

	return items, nil
}
