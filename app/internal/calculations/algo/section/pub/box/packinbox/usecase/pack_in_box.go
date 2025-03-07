package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"
	"github.com/mondegor/go-webcore/mrsender/decorator"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/dto"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/box/packinbox"
)

const (
	ModelNameBoxPackInBox = "public-api.Calculations.Algo.BoxRequest.PackInBox" // ModelNameBoxPackInBox - название сущности
)

type (
	// BoxPackInBox - comment struct.
	BoxPackInBox struct {
		algo         *packinbox.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewBoxPackInBox - создаёт объект BoxPackInBox.
func NewBoxPackInBox(algo *packinbox.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *BoxPackInBox {
	return &BoxPackInBox{
		algo:         algo,
		eventEmitter: decorator.NewSourceEmitter(eventEmitter, ModelNameBoxPackInBox),
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *BoxPackInBox) Calc(ctx context.Context, data dto.ParsedData) (model.BoxPackInBoxResponse, error) {
	result, err := uc.algo.Calc(data.Box, data.ProductHeap)
	if err != nil {
		return model.BoxPackInBoxResponse{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	var (
		fullBox model.BoxResponse
		restBox *model.BoxResponse
	)

	if !result.FullBox.IsEmpty() {
		fullBox = model.BoxResponse{
			Weight:              measure.Kilogram(mrlib.RoundFloat4(result.FullBox.TotalWeight())),
			Volume:              measure.Meter3(mrlib.RoundFloat8(result.FullBox.TotalVolume())),
			InnerVolume:         measure.Meter3(mrlib.RoundFloat8(result.FullBox.TotalInnerVolume())),
			ProductQuantity:     result.FullBox.Product.Quantity,
			ProductVolume:       measure.Meter3(mrlib.RoundFloat8(result.FullBox.Product.TotalVolume())),
			UnusedVolumePercent: mrlib.RoundFloat2(result.FullBox.UnusedVolumePercent()),
		}
	}

	if !result.RestBox.IsEmpty() {
		restBox = &model.BoxResponse{
			Weight:              measure.Kilogram(mrlib.RoundFloat4(result.RestBox.Weight())),
			Volume:              measure.Meter3(mrlib.RoundFloat8(result.RestBox.Box.Volume())),
			InnerVolume:         measure.Meter3(mrlib.RoundFloat8(result.RestBox.Box.InnerVolume())),
			ProductQuantity:     result.RestBox.Product.Quantity,
			ProductVolume:       measure.Meter3(mrlib.RoundFloat8(result.RestBox.Product.TotalVolume())),
			UnusedVolumePercent: mrlib.RoundFloat2(result.RestBox.UnusedVolumePercent()),
		}
	}

	uc.eventEmitter.Emit(ctx, "Calc", mrmsg.Data{"data": data})

	return model.BoxPackInBoxResponse{
		FullBox:          fullBox,
		RestBox:          restBox,
		BoxesQuantity:    result.BoxesQuantity(),
		BoxesWeight:      measure.Kilogram(mrlib.RoundFloat4(result.BoxesWeight())),
		ProductsVolume:   measure.Meter3(mrlib.RoundFloat8(result.ProductsVolume())),
		BoxesVolume:      measure.Meter3(mrlib.RoundFloat8(result.BoxesVolume())),
		BoxesInnerVolume: measure.Meter3(mrlib.RoundFloat8(result.BoxesInnerVolume())),
	}, nil
}
