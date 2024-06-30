package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

type (
	// RectImpositionUseCase - comment interface.
	RectImpositionUseCase interface {
		Calc(ctx context.Context, raw entity.RawData) (imposition.AlgoResult, error)
	}
)
