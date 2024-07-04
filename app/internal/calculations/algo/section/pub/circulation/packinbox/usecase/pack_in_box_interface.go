package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
)

type (
	// CirculationPackInBoxUseCase - comment interface.
	CirculationPackInBoxUseCase interface {
		CalcQuantity(ctx context.Context, raw entity.RawData) (entity.AlgoResult, error)
	}
)
