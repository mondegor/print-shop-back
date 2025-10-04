package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrargs"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlib/extmath"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/packinstack"
)

const (
	// ModelNameSheetPackInStack - название сущности.
	ModelNameSheetPackInStack = "public-api.Calculations.Algo.SheetRequest.PackInStack"
)

type (
	// SheetPackInStack - comment struct.
	SheetPackInStack struct {
		algo         *packinstack.AlgoSheet
		eventEmitter mrevent.Emitter
	}
)

// NewSheetPackInStack - создаёт объект SheetPackInStack.
func NewSheetPackInStack(algo *packinstack.AlgoSheet, eventEmitter mrevent.Emitter) *SheetPackInStack {
	return &SheetPackInStack{
		algo:         algo,
		eventEmitter: mrevent.NewSourceEmitter(eventEmitter, ModelNameSheetPackInStack),
	}
}

// Calc - comment method.
func (uc *SheetPackInStack) Calc(ctx context.Context, data dto.ParsedData) (model.SheetPackInStackResponse, error) {
	result, err := uc.algo.Calc(data.SheetHeap, data.QuantityInStack)
	if err != nil {
		return model.SheetPackInStackResponse{}, mr.ErrUseCaseIncorrectInputData.New(err)
	}

	var (
		fullBox model.ProductResponse
		restBox *model.ProductResponse
	)

	if !result.FullProduct.Empty() {
		fullBox = model.ProductResponse{
			Format: result.FullProduct.Format.Round(),
			Weight: measure.Kilogram(extmath.RoundFloat4(result.FullProduct.Product.Weight)),
			Volume: measure.Meter3(extmath.RoundFloat8(result.FullProduct.Product.Format.Volume())),
		}
	}

	if result.RestProduct.Weight != 0 {
		restBox = &model.ProductResponse{
			Format: result.RestProduct.Format.Round(),
			Weight: measure.Kilogram(extmath.RoundFloat4(result.RestProduct.Weight)),
			Volume: measure.Meter3(extmath.RoundFloat8(result.RestProduct.Format.Volume())),
		}
	}

	uc.eventEmitter.Emit(ctx, "Calc", mrargs.Group{"data": data})

	return model.SheetPackInStackResponse{
		FullProduct:   fullBox,
		RestProduct:   restBox,
		TotalQuantity: result.TotalQuantity(),
		TotalWeight:   measure.Kilogram(extmath.RoundFloat4(result.TotalWeight())),
		TotalVolume:   measure.Meter3(extmath.RoundFloat8(result.TotalVolume())),
	}, nil
}
