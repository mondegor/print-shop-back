package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrevent"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/util/conv"
	"github.com/mondegor/go-sysmess/util/xmath"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/controller/httpv1/model"
	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/box/packinbox/dto"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/algo/box/packinbox"
	"github.com/mondegor/print-shop-back/pkg/mrcalc/measure"
)

const (
	// ModelNameBoxPackInBox - название сущности.
	ModelNameBoxPackInBox = "public-api.Calculations.Algo.BoxRequest.PackInBox"
)

type (
	// BoxPackInBox - comment struct.
	BoxPackInBox struct {
		algo         *packinbox.Algo
		logger       mrlog.Logger
		eventEmitter mrevent.Emitter
	}
)

// NewBoxPackInBox - создаёт объект BoxPackInBox.
func NewBoxPackInBox(algo *packinbox.Algo, logger mrlog.Logger, eventEmitter mrevent.Emitter) *BoxPackInBox {
	return &BoxPackInBox{
		algo:         algo,
		logger:       logger,
		eventEmitter: mrevent.EmitterWithSource(eventEmitter, ModelNameBoxPackInBox),
	}
}

// Calc - comment method.
func (uc *BoxPackInBox) Calc(ctx context.Context, data dto.ParsedData) (model.BoxPackInBoxResponse, error) {
	result, err := uc.algo.Calc(ctx, data.Box, data.ProductHeap)
	if err != nil {
		return model.BoxPackInBoxResponse{}, errors.ErrUseCaseIncorrectInputData.New(err)
	}

	var (
		fullBox model.BoxResponse
		restBox *model.BoxResponse
	)

	if !result.FullBox.Empty() {
		fullBox = model.BoxResponse{
			Weight:              measure.Kilogram(xmath.RoundFloat4(result.FullBox.TotalWeight())),
			Volume:              measure.Meter3(xmath.RoundFloat8(result.FullBox.TotalVolume())),
			InnerVolume:         measure.Meter3(xmath.RoundFloat8(result.FullBox.TotalInnerVolume())),
			ProductQuantity:     result.FullBox.Product.Quantity,
			ProductVolume:       measure.Meter3(xmath.RoundFloat8(result.FullBox.Product.TotalVolume())),
			UnusedVolumePercent: xmath.RoundFloat2(result.FullBox.UnusedVolumePercent()),
		}
	}

	if !result.RestBox.Empty() {
		restBox = &model.BoxResponse{
			Weight:              measure.Kilogram(xmath.RoundFloat4(result.RestBox.Weight())),
			Volume:              measure.Meter3(xmath.RoundFloat8(result.RestBox.Box.Volume())),
			InnerVolume:         measure.Meter3(xmath.RoundFloat8(result.RestBox.Box.InnerVolume())),
			ProductQuantity:     result.RestBox.Product.Quantity,
			ProductVolume:       measure.Meter3(xmath.RoundFloat8(result.RestBox.Product.TotalVolume())),
			UnusedVolumePercent: xmath.RoundFloat2(result.RestBox.UnusedVolumePercent()),
		}
	}

	uc.eventEmitter.Emit(ctx, "Calc", conv.Group{"data": data})

	return model.BoxPackInBoxResponse{
		FullBox:          fullBox,
		RestBox:          restBox,
		BoxesQuantity:    result.BoxesQuantity(),
		BoxesWeight:      measure.Kilogram(xmath.RoundFloat4(result.BoxesWeight())),
		ProductsVolume:   measure.Meter3(xmath.RoundFloat8(result.ProductsVolume())),
		BoxesVolume:      measure.Meter3(xmath.RoundFloat8(result.BoxesVolume())),
		BoxesInnerVolume: measure.Meter3(xmath.RoundFloat8(result.BoxesInnerVolume())),
	}, nil
}
