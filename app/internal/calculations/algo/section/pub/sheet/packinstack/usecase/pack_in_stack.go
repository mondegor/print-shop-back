package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/sheet/packinstack/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/packinstack"
)

const (
	ModelNameSheetPackInStack = "public-api.Calculations.Algo.SheetRequest.PackInStack" // ModelNameSheetPackInStack - название сущности
)

type (
	// SheetPackInStack - comment struct.
	SheetPackInStack struct {
		algo         *packinstack.AlgoSheet
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewSheetPackInStack - создаёт объект SheetPackInStack.
func NewSheetPackInStack(algo *packinstack.AlgoSheet, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *SheetPackInStack {
	return &SheetPackInStack{
		algo:         algo,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, ModelNameSheetPackInStack),
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *SheetPackInStack) Calc(ctx context.Context, data dto.ParsedData) (model.SheetPackInStackResponse, error) {
	result, err := uc.algo.Calc(data.SheetHeap, data.QuantityInStack)
	if err != nil {
		return model.SheetPackInStackResponse{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	var (
		fullBox model.ProductResponse
		restBox *model.ProductResponse
	)

	if !result.FullProduct.IsEmpty() {
		fullBox = model.ProductResponse{
			Format: result.FullProduct.Format.Round(),
			Weight: measure.Kilogram(mrlib.RoundFloat4(result.FullProduct.Product.Weight)),
			Volume: measure.Meter3(mrlib.RoundFloat8(result.FullProduct.Product.Format.Volume())),
		}
	}

	if result.RestProduct.Weight != 0 {
		restBox = &model.ProductResponse{
			Format: result.RestProduct.Format.Round(),
			Weight: measure.Kilogram(mrlib.RoundFloat4(result.RestProduct.Weight)),
			Volume: measure.Meter3(mrlib.RoundFloat8(result.RestProduct.Format.Volume())),
		}
	}

	uc.eventEmitter.Emit(ctx, "Calc", mrmsg.Data{"data": data})

	return model.SheetPackInStackResponse{
		FullProduct:   fullBox,
		RestProduct:   restBox,
		TotalQuantity: result.TotalQuantity(),
		TotalWeight:   measure.Kilogram(mrlib.RoundFloat4(result.TotalWeight())),
		TotalVolume:   measure.Meter3(mrlib.RoundFloat8(result.TotalVolume())),
	}, nil
}
