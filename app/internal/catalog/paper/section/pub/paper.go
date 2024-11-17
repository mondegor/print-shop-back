package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/catalog/paper/section/pub/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// PaperUseCase - comment interface.
	PaperUseCase interface {
		GetList(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
		GetTypeList(ctx context.Context) ([]uint64, error)
		GetColorList(ctx context.Context) ([]uint64, error)
		GetDensityList(ctx context.Context) ([]measure.KilogramPerMeter2, error)
		GetFactureList(ctx context.Context) ([]uint64, error)
	}

	// PaperStorage - comment interface.
	PaperStorage interface {
		Fetch(ctx context.Context, params entity.PaperParams) ([]entity.Paper, error)
		FetchTypeIDs(ctx context.Context) ([]uint64, error)
		FetchColorIDs(ctx context.Context) ([]uint64, error)
		FetchDensities(ctx context.Context) ([]measure.KilogramPerMeter2, error)
		FetchFactureIDs(ctx context.Context) ([]uint64, error)
	}
)
