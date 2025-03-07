package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/cutting/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/cutting"
)

const (
	ModelNameSheetCutting = "public-api.Calculations.Algo.Sheet.Cutting" // ModelNameSheetCutting - название сущности
)

type (
	// SheetCutting - comment struct.
	SheetCutting struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewSheetCutting - создаёт объект SheetCutting.
func NewSheetCutting(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *SheetCutting {
	return &SheetCutting{
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, ModelNameSheetCutting),
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *SheetCutting) CalcQuantity(ctx context.Context, data dto.ParsedData) (model.SheetCuttingQuantityResult, error) {
	result := cutting.AlgoQuantity(data.Fragments, data.DistanceFormat)

	uc.eventEmitter.Emit(ctx, "CalcQuantity", mrmsg.Data{"data": data})

	return model.SheetCuttingQuantityResult{
		Quantity: result,
	}, nil
}
