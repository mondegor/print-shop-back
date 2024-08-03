package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
)

type (
	// RectInsideOutsideUseCase - comment interface.
	RectInsideOutsideUseCase interface {
		CalcQuantity(ctx context.Context, data entity.ParsedData) (entity.AlgoQuantityResult, error)
		CalcMax(ctx context.Context, data entity.ParsedData) (entity.AlgoMaxResult, error)
	}
)
