package usecase

import (
	"context"
	"math"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/util/conv"
	"github.com/mondegor/go-sysmess/util/xmath"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/dto"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/measure"
)

const (
	// ModelNameSheetImposition - название сущности.
	ModelNameSheetImposition = "public-api.Calculations.Algo.Sheet.Imposition"
)

type (
	// SheetImposition - comment struct.
	SheetImposition struct {
		algo *imposition.Algo
		// logger       mrlog.Logger
		eventEmitter mrevent.Emitter
	}
)

// NewSheetImposition - создаёт объект SheetImposition.
func NewSheetImposition(algo *imposition.Algo, eventEmitter mrevent.Emitter) *SheetImposition {
	return &SheetImposition{
		algo: algo,
		// logger:       logger,
		eventEmitter: mrevent.EmitterWithSource(eventEmitter, ModelNameSheetImposition),
	}
}

// Calc - comment method.
func (uc *SheetImposition) Calc(ctx context.Context, data dto.ParsedData) (model.SheetImpositionResponse, error) {
	result, err := uc.algo.Calc(ctx, data.Element, data.Distance, data.Out, data.Opts)
	if err != nil {
		return model.SheetImpositionResponse{}, errors.ErrIncorrectInputData.New(err)
	}

	uc.eventEmitter.Emit(ctx, "Calc", conv.Group{"data": data})

	if result.Fragments.TotalQuantity() == 0 {
		return model.SheetImpositionResponse{}, nil
	}

	return uc.createResult(result), nil
}

// CalcVariants - comment method.
func (uc *SheetImposition) CalcVariants(ctx context.Context, data dto.ParsedData) (model.SheetImpositionVariantsResponse, error) {
	if almostEqual(data.Element.Width, data.Element.Height) && almostEqual(data.Distance.Width, data.Distance.Height) {
		result, err := uc.algo.Calc(ctx, data.Element, data.Distance, data.Out, data.Opts)
		if err != nil {
			return nil, errors.ErrIncorrectInputData.New(err)
		}

		uc.eventEmitter.Emit(ctx, "CalcVariants", conv.Group{"data": data})

		return model.SheetImpositionVariantsResponse{
			uc.createResult(result),
		}, nil
	}

	countVariants := 0

	resultV1, errV1 := uc.algo.Calc(ctx, data.Element, data.Distance, data.Out, data.Opts)
	if errV1 == nil {
		countVariants++
	}

	resultV2, errV2 := uc.algo.Calc(ctx, data.Element.Rotate90(), data.Distance.Rotate90(), data.Out, data.Opts)
	if errV2 == nil {
		countVariants++
	}

	if countVariants == 0 {
		return nil, errors.ErrIncorrectInputData.New(errV1)
	}

	results := make(model.SheetImpositionVariantsResponse, 0, countVariants)

	if errV1 == nil {
		results = append(results, uc.createResult(resultV1))
	}

	if errV2 == nil {
		results = append(results, uc.createResult(resultV2))
	}

	uc.eventEmitter.Emit(ctx, "CalcVariants", conv.Group{"data": data})

	return results, nil
}

func (uc *SheetImposition) createResult(item imposition.Output) model.SheetImpositionResponse {
	return model.SheetImpositionResponse{
		ContainerFormat:  item.ContainerFormat.Round(),
		FragmentDistance: measure.Meter(xmath.RoundFloat4(item.Fragments.FragmentDistance())),
		Fragments:        item.Fragments.Round(),
		TotalElements:    item.Fragments.TotalQuantity(),
		Garbage:          measure.Meter2(xmath.RoundFloat8(item.RestArea)),
		AllowRotation:    item.AllowRotation,
		UseMirror:        item.UseMirror,
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}
