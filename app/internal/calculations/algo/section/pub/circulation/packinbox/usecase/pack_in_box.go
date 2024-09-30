package usecase

import (
	"context"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/packinbox"
)

type (
	// CirculationPackInBox - comment struct.
	CirculationPackInBox struct {
		algo         *packinbox.Algo
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UseCaseErrorWrapper
	}
)

// NewCirculationPackInBox - создаёт объект CirculationPackInBox.
func NewCirculationPackInBox(algo *packinbox.Algo, eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UseCaseErrorWrapper) *CirculationPackInBox {
	return &CirculationPackInBox{
		algo:         algo,
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *CirculationPackInBox) Calc(ctx context.Context, data entity.ParsedData) (entity.AlgoResult, error) {
	result, err := uc.algo.Calc(data.Box, data.Product)
	if err != nil {
		return entity.AlgoResult{}, mrcore.ErrUseCaseIncorrectInputData.Wrap(err, "data", data)
	}

	var (
		fullBox entity.BoxResult
		restBox entity.BoxResult
	)

	if result.FullBox.ProductQuantity > 0 {
		fullBox = entity.BoxResult{
			Weight:              mrlib.RoundFloat4(result.FullBox.Weight),
			Volume:              mrlib.RoundFloat8(result.FullBox.Volume),
			InnerVolume:         mrlib.RoundFloat8(result.FullBox.InnerVolume),
			ProductQuantity:     result.FullBox.ProductQuantity,
			ProductVolume:       mrlib.RoundFloat8(result.FullBox.ProductVolume),
			UnusedVolumePercent: mrlib.RoundFloat2(result.FullBox.UnusedVolumePercent),
		}
	}

	if result.RestBox.ProductQuantity > 0 {
		restBox = entity.BoxResult{
			Weight:              mrlib.RoundFloat4(result.RestBox.Weight),
			Volume:              mrlib.RoundFloat8(result.RestBox.Volume),
			InnerVolume:         mrlib.RoundFloat8(result.RestBox.InnerVolume),
			ProductQuantity:     result.RestBox.ProductQuantity,
			ProductVolume:       mrlib.RoundFloat8(result.RestBox.ProductVolume),
			UnusedVolumePercent: mrlib.RoundFloat2(result.RestBox.UnusedVolumePercent),
		}
	}

	uc.emitEvent(ctx, "Calc", mrmsg.Data{"data": data})

	return entity.AlgoResult{
		FullBox:          fullBox,
		RestBox:          restBox,
		BoxesQuantity:    result.BoxesQuantity,
		BoxesWeight:      mrlib.RoundFloat4(result.BoxesWeight),
		ProductsVolume:   mrlib.RoundFloat8(result.ProductsVolume),
		BoxesVolume:      mrlib.RoundFloat8(result.BoxesVolume),
		BoxesInnerVolume: mrlib.RoundFloat8(result.BoxesInnerVolume),
	}, nil
}

func (uc *CirculationPackInBox) emitEvent(ctx context.Context, eventName string, data mrmsg.Data) {
	uc.eventEmitter.EmitWithSource(
		ctx,
		eventName,
		entity.ModelNameCirculationPackInBox,
		data,
	)
}
