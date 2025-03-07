package usecase

import (
	"context"
	"math"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/imposition/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
)

const (
	ModelNameSheetImposition = "public-api.Calculations.Algo.Sheet.Imposition" // ModelNameSheetImposition - название сущности
)

type (
	// SheetImposition - comment struct.
	SheetImposition struct {
		algo         *imposition.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewSheetImposition - создаёт объект SheetImposition.
func NewSheetImposition(algo *imposition.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *SheetImposition {
	return &SheetImposition{
		algo:         algo,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, ModelNameSheetImposition),
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *SheetImposition) Calc(ctx context.Context, data dto.ParsedData) (model.SheetImpositionResponse, error) {
	result, err := uc.algo.Calc(data.Element, data.Distance, data.Out, data.Opts)
	if err != nil {
		return model.SheetImpositionResponse{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.eventEmitter.Emit(ctx, "Calc", mrmsg.Data{"data": data})

	if result.Fragments.TotalQuantity() == 0 {
		return model.SheetImpositionResponse{}, nil
	}

	return uc.createResult(result), nil
}

// CalcVariants - comment method.
func (uc *SheetImposition) CalcVariants(ctx context.Context, data dto.ParsedData) (model.SheetImpositionVariantsResponse, error) {
	if almostEqual(data.Element.Width, data.Element.Height) && almostEqual(data.Distance.Width, data.Distance.Height) {
		result, err := uc.algo.Calc(data.Element, data.Distance, data.Out, data.Opts)
		if err != nil {
			return nil, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
		}

		uc.eventEmitter.Emit(ctx, "CalcVariants", mrmsg.Data{"data": data})

		return model.SheetImpositionVariantsResponse{
			uc.createResult(result),
		}, nil
	}

	var countVariants int

	resultV1, errV1 := uc.algo.Calc(data.Element, data.Distance, data.Out, data.Opts)
	if errV1 == nil {
		countVariants++
	}

	resultV2, errV2 := uc.algo.Calc(data.Element.Rotate90(), data.Distance.Rotate90(), data.Out, data.Opts)
	if errV2 == nil {
		countVariants++
	}

	if countVariants == 0 {
		return nil, mrcore.ErrUseCaseIncorrectInputData.Wrap(errV1, "data", data)
	}

	results := make(model.SheetImpositionVariantsResponse, 0, countVariants)

	if errV1 == nil {
		results = append(results, uc.createResult(resultV1))
	}

	if errV2 == nil {
		results = append(results, uc.createResult(resultV2))
	}

	uc.eventEmitter.Emit(ctx, "CalcVariants", mrmsg.Data{"data": data})

	return results, nil
}

func (uc *SheetImposition) createResult(item imposition.Output) model.SheetImpositionResponse {
	return model.SheetImpositionResponse{
		ContainerFormat:  item.ContainerFormat.Round(),
		FragmentDistance: measure.Meter(mrlib.RoundFloat4(item.Fragments.FragmentDistance())),
		Fragments:        item.Fragments.Round(),
		TotalElements:    item.Fragments.TotalQuantity(),
		Garbage:          measure.Meter2(mrlib.RoundFloat8(item.RestArea)),
		AllowRotation:    item.AllowRotation,
		UseMirror:        item.UseMirror,
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-9
}
