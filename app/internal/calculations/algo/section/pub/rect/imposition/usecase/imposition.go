package usecase

import (
	"context"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
)

type (
	// RectImposition - comment struct.
	RectImposition struct {
		algo         *imposition.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewRectImposition - создаёт объект RectImposition.
func NewRectImposition(algo *imposition.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *RectImposition {
	return &RectImposition{
		algo:         algo,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *RectImposition) Calc(ctx context.Context, data entity.ParsedData) (entity.AlgoResult, error) {
	result, err := uc.algo.Calc(data.Item, data.Out, data.Opts)
	if err != nil {
		return entity.AlgoResult{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.emitEvent(ctx, "Calc", mrmsg.Data{"data": data})

	if result.Total == 0 {
		return entity.AlgoResult{}, nil
	}

	return entity.AlgoResult{
		Layout: rect.Format{
			Width:  mrlib.RoundFloat4(result.Layout.Width),
			Height: mrlib.RoundFloat4(result.Layout.Height),
		},
		Fragments: result.Fragments,
		Total:     result.Total,
		Garbage:   mrlib.RoundFloat8(result.RestArea),
	}, nil
}

func (uc *RectImposition) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameRectImposition,
		data,
	)
}
