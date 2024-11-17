package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/cutting"
)

type (
	// RectCutting - comment struct.
	RectCutting struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewRectCutting - создаёт объект RectCutting.
func NewRectCutting(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *RectCutting {
	return &RectCutting{
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameRectCutting),
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *RectCutting) CalcQuantity(ctx context.Context, data entity.ParsedData) (entity.QuantityResult, error) {
	result := cutting.AlgoQuantity(data.Fragments, data.DistanceFormat)

	uc.eventEmitter.Emit(ctx, "CalcQuantity", mrmsg.Data{"data": data})

	return entity.QuantityResult{
		Quantity: result,
	}, nil
}
