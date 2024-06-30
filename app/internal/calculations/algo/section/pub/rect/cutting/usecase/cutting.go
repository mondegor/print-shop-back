package usecase

import (
	"context"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/cutting/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/base"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
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
func (uc *RectCutting) CalcQuantity(ctx context.Context, raw entity.RawData) (entity.AlgoQuantityResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.AlgoQuantityResult{}, err
	}

	result := cutting.AlgoQuantity(parsedData.Fragments, parsedData.DistanceFormat)

	uc.emitEvent(ctx, "CalcQuantity", mrmsg.Data{"raw": parsedData})

	return entity.AlgoQuantityResult{
		Quantity: result,
	}, nil
}

func (uc *RectCutting) parse(data entity.RawData) (entity.ParsedData, error) {
	fragments := make([]base.Fragment, 0, len(data.Fragments))

	for _, str := range data.Fragments {
		fragment, err := base.ParseFragment(str)
		if err != nil {
			return entity.ParsedData{}, err // TODO: itemInFormat error
		}

		fragments = append(fragments, fragment)
	}

	distanceFormat, err := rect.ParseFormat(data.DistanceFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: distanceFormat error
	}

	return entity.ParsedData{
		Fragments:      fragments,
		DistanceFormat: distanceFormat,
	}, nil
}

func (uc *RectCutting) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCutting,
		data,
	)
}
