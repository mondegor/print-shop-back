package pub

import (
	"context"

	boxpackinboxmodel "github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1/model"
	boxpackinboxdto "github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/dto"
)

type (
	// BoxPackInBoxUseCase - comment interface.
	BoxPackInBoxUseCase interface {
		Calc(ctx context.Context, data boxpackinboxdto.ParsedData) (boxpackinboxmodel.BoxPackInBoxResponse, error)
	}
)
