package pub

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
)

type (
	// CirculationPackInBoxUseCase - comment interface.
	CirculationPackInBoxUseCase interface {
		Calc(ctx context.Context, data entity.ParsedData) (entity.AlgoResult, error)
	}
)
