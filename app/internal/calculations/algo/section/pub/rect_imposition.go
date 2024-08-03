package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
)

type (
	// RectImpositionUseCase - comment interface.
	RectImpositionUseCase interface {
		Calc(ctx context.Context, data entity.ParsedData) (entity.AlgoResult, error)
	}
)
