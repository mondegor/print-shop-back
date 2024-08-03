package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/cutting"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// RectCutting - comment struct.
	RectCutting struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewRectCutting - создаёт объект RectCutting.
func NewRectCutting(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *RectCutting {
	return &RectCutting{
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *RectCutting) CalcQuantity(ctx context.Context, data entity.ParsedData) (entity.QuantityResult, error) {
	result := cutting.AlgoQuantity(data.Fragments, data.DistanceFormat)

	uc.emitEvent(ctx, "CalcQuantity", mrmsg.Data{"data": data})

	return entity.QuantityResult{
		Quantity: result,
	}, nil
}

func (uc *RectCutting) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameRectCutting,
		data,
	)
}
