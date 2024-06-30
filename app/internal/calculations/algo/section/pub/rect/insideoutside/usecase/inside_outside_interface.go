package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
)

type (
	// RectInsideOutsideUseCase - comment interface.
	RectInsideOutsideUseCase interface {
		CalcQuantity(ctx context.Context, raw entity.RawData) (entity.AlgoQuantityResult, error)
		CalcMax(ctx context.Context, raw entity.RawData) (entity.AlgoMaxResult, error)
	}
)
