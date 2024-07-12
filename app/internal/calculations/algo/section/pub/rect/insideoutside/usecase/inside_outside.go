package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/insideoutside/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/insideoutside"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// RectInsideOutside - comment struct.
	RectInsideOutside struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewRectInsideOutside - создаёт объект RectInsideOutside.
func NewRectInsideOutside(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *RectInsideOutside {
	return &RectInsideOutside{
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *RectInsideOutside) CalcQuantity(ctx context.Context, raw entity.RawData) (entity.AlgoQuantityResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.AlgoQuantityResult{}, err
	}

	result, err := insideoutside.AlgoQuantity(parsedData.In, parsedData.Out)
	if err != nil {
		return entity.AlgoQuantityResult{}, err
	}

	uc.emitEvent(ctx, "CalcQuantity", mrmsg.Data{"raw": parsedData})

	return entity.AlgoQuantityResult{
		Fragment: result,
		Total:    result.Total(),
	}, nil
}

// CalcMax - comment method.
func (uc *RectInsideOutside) CalcMax(ctx context.Context, raw entity.RawData) (entity.AlgoMaxResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.AlgoMaxResult{}, err
	}

	result, err := insideoutside.AlgoMax(parsedData.In, parsedData.Out)
	if err != nil {
		return entity.AlgoMaxResult{}, err
	}

	uc.emitEvent(ctx, "CalcMax", mrmsg.Data{"raw": parsedData})

	return entity.AlgoMaxResult{
		Fragments: result,
		Total:     result.Total(),
	}, nil
}

func (uc *RectInsideOutside) parse(data entity.RawData) (entity.ParsedData, error) {
	inFormat, err := rect.ParseFormat(data.InFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemInFormat error
	}

	outFormat, err := rect.ParseFormat(data.OutFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemOutFormat error
	}

	return entity.ParsedData{
		In: rect.Format{
			Width:  inFormat.Width * measure.OneThousandth,  // mm -> m
			Height: inFormat.Height * measure.OneThousandth, // mm -> m
		},
		Out: rect.Format{
			Width:  outFormat.Width * measure.OneThousandth,  // mm -> m
			Height: outFormat.Height * measure.OneThousandth, // mm -> m
		},
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
