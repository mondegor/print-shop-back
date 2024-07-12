package usecase

import (
	"context"
	"errors"

	"github.com/mondegor/go-webcore/mrlib"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
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
func (uc *RectImposition) Calc(ctx context.Context, raw entity.RawData) (entity.Result, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.Result{}, err
	}

	result, err := uc.algo.Calc(parsedData.Item, parsedData.Out, parsedData.Opts)
	if err != nil {
		return entity.Result{}, err
	}

	if len(result.Fragments) == 0 {
		return entity.Result{}, errors.New("error") // TODO: result NULL
	}

	uc.emitEvent(ctx, "Calc", mrmsg.Data{"raw": parsedData})

	return entity.Result{
		Layout: rect.Format{
			Width:  mrlib.RoundFloat4(result.Layout.Width),
			Height: mrlib.RoundFloat4(result.Layout.Height),
		},
		Fragments: result.Fragments,
		Total:     result.Total,
		Garbage:   mrlib.RoundFloat8(result.RestArea),
	}, nil
}

func (uc *RectImposition) parse(data entity.RawData) (entity.ParsedData, error) {
	itemFormat, err := rect.ParseFormat(data.ItemFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: itemFormat error
	}

	itemDistance := rect.Format{} // optional

	if data.ItemDistance != "" {
		itemDistance, err = rect.ParseFormat(data.ItemDistance)
		if err != nil {
			return entity.ParsedData{}, err // TODO: itemDistance error
		}
	}

	outFormat, err := rect.ParseFormat(data.OutFormat)
	if err != nil {
		return entity.ParsedData{}, err // TODO: outFormat error
	}

	return entity.ParsedData{
		Item: rect.Item{
			Format: rect.Format{
				Width:  itemFormat.Width * measure.OneThousandth,  // mm -> m
				Height: itemFormat.Height * measure.OneThousandth, // mm -> m
			},
			Distance: rect.Format{
				Width:  itemDistance.Width * measure.OneThousandth,  // mm -> m
				Height: itemDistance.Height * measure.OneThousandth, // mm -> m
			},
		},
		Out: rect.Format{
			Width:  outFormat.Width * measure.OneThousandth,  // mm -> m
			Height: outFormat.Height * measure.OneThousandth, // mm -> m
		},
		Opts: imposition.Options{
			AllowRotation: data.AllowRotation,
			UseMirror:     data.UseMirror,
		},
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
