package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// PaperUseCase - comment interface.
	PaperUseCase interface {
		GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
		GetTypeList(ctx context.Context) ([]mrtype.KeyInt32, error)
		GetColorList(ctx context.Context) ([]mrtype.KeyInt32, error)
		GetDensityList(ctx context.Context) ([]measure.GramsPerMeter2, error)
		GetFactureList(ctx context.Context) ([]mrtype.KeyInt32, error)
	}

	// PaperStorage - comment interface.
	PaperStorage interface {
		Fetch(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
		FetchTypeIDs(ctx context.Context) ([]mrtype.KeyInt32, error)
		FetchColorIDs(ctx context.Context) ([]mrtype.KeyInt32, error)
		FetchDensities(ctx context.Context) ([]measure.GramsPerMeter2, error)
		FetchFactureIDs(ctx context.Context) ([]mrtype.KeyInt32, error)
	}
)
