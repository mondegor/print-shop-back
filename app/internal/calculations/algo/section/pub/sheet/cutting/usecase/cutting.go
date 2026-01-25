package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/util/conv"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/dto"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/sheet/cutting"
)

const (
	// ModelNameSheetCutting - название сущности.
	ModelNameSheetCutting = "public-api.Calculations.Algo.Sheet.Cutting"
)

type (
	// SheetCutting - comment struct.
	SheetCutting struct {
		eventEmitter mrevent.Emitter
	}
)

// NewSheetCutting - создаёт объект SheetCutting.
func NewSheetCutting(eventEmitter mrevent.Emitter) *SheetCutting {
	return &SheetCutting{
		eventEmitter: mrevent.EmitterWithSource(eventEmitter, ModelNameSheetCutting),
	}
}

// CalcQuantity - comment method.
func (uc *SheetCutting) CalcQuantity(ctx context.Context, data dto.ParsedData) (model.SheetCuttingQuantityResult, error) {
	result := cutting.AlgoQuantity(data.Fragments, data.DistanceFormat)

	uc.eventEmitter.Emit(ctx, "CalcQuantity", conv.Group{"data": data})

	return model.SheetCuttingQuantityResult{
		Quantity: result,
	}, nil
}
