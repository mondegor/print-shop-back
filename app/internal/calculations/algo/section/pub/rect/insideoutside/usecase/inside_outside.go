package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"
)

type (
	// RectInsideOutside - comment struct.
	RectInsideOutside struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewRectInsideOutside - создаёт объект RectInsideOutside.
func NewRectInsideOutside(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *RectInsideOutside {
	return &RectInsideOutside{
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *RectInsideOutside) CalcQuantity(ctx context.Context, data entity.ParsedData) (entity.AlgoQuantityResult, error) {
	result, err := insideoutside.AlgoQuantity(data.In, data.Out)
	if err != nil {
		return entity.AlgoQuantityResult{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.emitEvent(ctx, "CalcQuantity", mrmsg.Data{"data": data})

	return entity.AlgoQuantityResult{
		Fragment: result,
		Total:    result.Total(),
	}, nil
}

// CalcMax - comment method.
func (uc *RectInsideOutside) CalcMax(ctx context.Context, data entity.ParsedData) (entity.AlgoMaxResult, error) {
	result, err := insideoutside.AlgoMax(data.In, data.Out)
	if err != nil {
		return entity.AlgoMaxResult{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.emitEvent(ctx, "CalcMax", mrmsg.Data{"data": data})

	return entity.AlgoMaxResult{
		Fragments: result,
		Total:     result.Total(),
	}, nil
}

func (uc *RectInsideOutside) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameRectInsideOutside,
		data,
	)
}
