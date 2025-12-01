package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/insideoutside/dto"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/sheet/insideoutside"
)

const (
	// ModelNameSheetInsideOutside - название сущности.
	ModelNameSheetInsideOutside = "public-api.Calculations.Algo.Sheet.InsideOutside"
)

type (
	// SheetInsideOutside - comment struct.
	SheetInsideOutside struct {
		eventEmitter mrevent.Emitter
	}
)

// NewSheetInsideOutside - создаёт объект SheetInsideOutside.
func NewSheetInsideOutside(eventEmitter mrevent.Emitter) *SheetInsideOutside {
	return &SheetInsideOutside{
		eventEmitter: mrevent.NewSourceEmitter(eventEmitter, ModelNameSheetInsideOutside),
	}
}

// CalcQuantity - comment method.
func (uc *SheetInsideOutside) CalcQuantity(ctx context.Context, data dto.ParsedData) (model.SheetInsideOutsideQuantityResponse, error) {
	result, err := insideoutside.AlgoQuantity(data.In, data.Out)
	if err != nil {
		return model.SheetInsideOutsideQuantityResponse{}, mr.ErrUseCaseIncorrectInputData.New(err)
	}

	uc.eventEmitter.Emit(ctx, "CalcQuantity", mrargs.Group{"data": data})

	return model.SheetInsideOutsideQuantityResponse{
		Layout: result,
		Total:  result.Quantity(),
	}, nil
}

// CalcMax - comment method.
func (uc *SheetInsideOutside) CalcMax(ctx context.Context, data dto.ParsedData) (model.SheetInsideOutsideMaxResponse, error) {
	result, err := insideoutside.AlgoMax(data.In, data.Out)
	if err != nil {
		return model.SheetInsideOutsideMaxResponse{}, mr.ErrUseCaseIncorrectInputData.New(err)
	}

	uc.eventEmitter.Emit(ctx, "CalcMax", mrargs.Group{"data": data})

	return model.SheetInsideOutsideMaxResponse{
		Fragments: result,
		Total:     result.TotalQuantity(),
	}, nil
}
