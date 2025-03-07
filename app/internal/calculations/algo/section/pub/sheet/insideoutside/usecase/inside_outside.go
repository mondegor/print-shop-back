package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/insideoutside"
)

const (
	ModelNameSheetInsideOutside = "public-api.Calculations.Algo.Sheet.InsideOutside" // ModelNameSheetInsideOutside - название сущности
)

type (
	// SheetInsideOutside - comment struct.
	SheetInsideOutside struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewSheetInsideOutside - создаёт объект SheetInsideOutside.
func NewSheetInsideOutside(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *SheetInsideOutside {
	return &SheetInsideOutside{
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, ModelNameSheetInsideOutside),
		errorWrapper: errorWrapper,
	}
}

// CalcQuantity - comment method.
func (uc *SheetInsideOutside) CalcQuantity(ctx context.Context, data dto.ParsedData) (model.SheetInsideOutsideQuantityResponse, error) {
	result, err := insideoutside.AlgoQuantity(data.In, data.Out)
	if err != nil {
		return model.SheetInsideOutsideQuantityResponse{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.eventEmitter.Emit(ctx, "CalcQuantity", mrmsg.Data{"data": data})

	return model.SheetInsideOutsideQuantityResponse{
		Layout: result,
		Total:  result.Quantity(),
	}, nil
}

// CalcMax - comment method.
func (uc *SheetInsideOutside) CalcMax(ctx context.Context, data dto.ParsedData) (model.SheetInsideOutsideMaxResponse, error) {
	result, err := insideoutside.AlgoMax(data.In, data.Out)
	if err != nil {
		return model.SheetInsideOutsideMaxResponse{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	uc.eventEmitter.Emit(ctx, "CalcMax", mrmsg.Data{"data": data})

	return model.SheetInsideOutsideMaxResponse{
		Fragments: result,
		Total:     result.TotalQuantity(),
	}, nil
}
