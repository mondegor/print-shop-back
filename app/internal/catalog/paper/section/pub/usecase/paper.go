package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub"
	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// Paper - comment struct.
	Paper struct {
		storage      pub.PaperStorage
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewPaper - создаёт объект Paper.
func NewPaper(storage pub.PaperStorage, errorWrapper mrcore.UseCaseErrorWrapper) *Paper {
	return &Paper{
		storage:      storage,
		errorWrapper: errorWrapper,
	}
}

// GetList - comment method.
func (uc *Paper) GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error) {
	items, err := uc.storage.Fetch(ctx, params)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}

// GetTypeList - comment method.
func (uc *Paper) GetTypeList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchTypeIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}

// GetColorList - comment method.
func (uc *Paper) GetColorList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchColorIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}

// GetDensityList - comment method.
func (uc *Paper) GetDensityList(ctx context.Context) ([]measure.KilogramPerMeter2, error) {
	items, err := uc.storage.FetchDensities(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}

// GetFactureList - comment method.
func (uc *Paper) GetFactureList(ctx context.Context) ([]uint64, error) {
	items, err := uc.storage.FetchFactureIDs(ctx)
	if err != nil {
		return nil, uc.errorWrapper.WrapErrorFailed(err, entity.ModelNamePaper)
	}

	return items, nil
}
