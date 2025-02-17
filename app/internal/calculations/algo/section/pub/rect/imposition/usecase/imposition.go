package usecase

import (
	"context"
	"math"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/rect/imposition/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

type (
	// RectImposition - comment struct.
	RectImposition struct {
		algo         *imposition.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}

	variant struct {
		itm *rect.Item
		res *imposition.AlgoResult
	}
)

// NewRectImposition - создаёт объект RectImposition.
func NewRectImposition(algo *imposition.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *RectImposition {
	return &RectImposition{
		algo:         algo,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, entity.ModelNameRectImposition),
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *RectImposition) Calc(ctx context.Context, data entity.ParsedData) (entity.AlgoResult, error) {
	result, err := uc.algo.Calc(data.Item, data.Out, data.Opts)
	if err != nil {
		return entity.AlgoResult{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.eventEmitter.Emit(ctx, "Calc", mrmsg.Data{"data": data})

	if result.Total == 0 {
		return entity.AlgoResult{}, nil
	}

	return uc.createResult(result), nil
}

func (uc *RectImposition) CalcVariants(ctx context.Context, data entity.ParsedData) ([]entity.AlgoVariantResult, error) {
	if almostEqual(data.Item.Format.Width, data.Item.Format.Height) && almostEqual(data.Item.Distance.Width, data.Item.Distance.Height) {
		result, err := uc.algo.Calc(data.Item, data.Out, data.Opts)
		if err != nil {
			return nil, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
		}

		uc.eventEmitter.Emit(ctx, "CalcVariants", mrmsg.Data{"data": data})

		return []entity.AlgoVariantResult{
			{
				Name:   data.Item.Format.Sum(data.Item.Distance).String() + " -> " + data.Out.String(),
				Result: uc.createResult(result),
			},
		}, nil
	}

	itemV2 := rect.Item{
		Format:   data.Item.Format.Change(),
		Distance: data.Item.Distance.Change(),
	}

	resultV1, errV1 := uc.algo.Calc(data.Item, data.Out, data.Opts)
	resultV2, errV2 := uc.algo.Calc(itemV2, data.Out, data.Opts)

	if errV1 != nil && errV2 != nil {
		return nil, mrcore.ErrUseCaseIncorrectInputData.Wrap(errV1, "data", data)
	}

	uc.eventEmitter.Emit(ctx, "CalcVariants", mrmsg.Data{"data": data})

	variants := make([]variant, 0, 2)

	if errV1 == nil && resultV1.Total > 0 {
		variants = append(
			variants,
			variant{
				itm: &data.Item,
				res: &resultV1,
			},
		)
	}

	if errV2 == nil && resultV2.Total > 0 {
		variants = append(
			variants,
			variant{
				itm: &itemV2,
				res: &resultV2,
			},
		)
	}

	if len(variants) == 2 && variants[0].res.RestArea > variants[1].res.RestArea {
		variants[0], variants[1] = variants[1], variants[0]
	}

	results := make([]entity.AlgoVariantResult, len(variants))

	results[0] = entity.AlgoVariantResult{
		Name:   variants[0].itm.Format.Sum(variants[0].itm.Distance).String() + " -> " + data.Out.String(),
		Result: uc.createResult(*variants[0].res),
	}

	if len(variants) == 2 {
		results[1] = entity.AlgoVariantResult{
			Name:   variants[1].itm.Format.Sum(variants[1].itm.Distance).String() + " -> " + data.Out.String(),
			Result: uc.createResult(*variants[1].res),
		}
	}

	return results, nil
}

func (uc *RectImposition) createResult(item imposition.AlgoResult) entity.AlgoResult {
	return entity.AlgoResult{
		Layout: rect.Format{
			Width:  mrlib.RoundFloat4(item.Layout.Width),
			Height: mrlib.RoundFloat4(item.Layout.Height),
		},
		Item:      item.Item,
		Distance:  item.Distance,
		Fragments: item.Fragments,
		Total:     item.Total,
		Garbage:   mrlib.RoundFloat8(item.RestArea),
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}
