package usecase

import (
	"context"
	"errors"

	"github.com/mondegor/go-sysmess/mrmsg"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrsender"

	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/parallelepiped"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"

	"github.com/mondegor/print-shop-back/internal/calculations/algo/section/pub/circulation/packinbox/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
)

type (
	// CirculationPackInBox - comment struct.
	CirculationPackInBox struct {
		eventEmitter mrsender.EventEmitter
		errorWrapper mrcore.UsecaseErrorWrapper
	}
)

// NewCirculationPackInBox - создаёт объект CirculationPackInBox.
func NewCirculationPackInBox(eventEmitter mrsender.EventEmitter, errorWrapper mrcore.UsecaseErrorWrapper) *CirculationPackInBox {
	return &CirculationPackInBox{
		eventEmitter: eventEmitter,
		errorWrapper: errorWrapper,
	}
}

// Calc - comment method.
func (uc *CirculationPackInBox) Calc(ctx context.Context, raw entity.RawData) (entity.AlgoResult, error) {
	parsedData, err := uc.parse(raw)
	if err != nil {
		return entity.AlgoResult{}, err
	}

	imp := imposition.New(mrlog.Ctx(ctx))

	boxBottomFormat := parsedData.Box.Format.BottomFormat()
	boxBottomFormat = boxBottomFormat.Diff(parsedData.Box.Margins.BottomFormat())

	// вычитается толщина стенки
	const twoSides = 2
	boxBottomFormat = boxBottomFormat.Diff(
		rect.Format{
			Width:  parsedData.Box.Thickness * twoSides,
			Height: parsedData.Box.Thickness * twoSides,
		},
	)

	impResult, err := imp.Calc(
		rect.Item{
			Format: parsedData.Product.Format,
		},
		boxBottomFormat,
		imposition.Options{
			AllowRotation: true,
			UseMirror:     false,
		},
	)
	if err != nil {
		return entity.AlgoResult{}, err // TODO: обернуть
	}

	if impResult.Total == 0 {
		return entity.AlgoResult{}, errors.New("total = 0") // коробка не подходит
	}

	// :WARNING: parsedData.Product.Thickness > 0 && parsedData.Product.WeightM2 > 0

	var (
		box     entity.BoxResult
		lastBox entity.BoxResult
	)

	// максимальное кол-во изделий в стопке (отбрасывается дробная часть)
	maxProductQuantityInStack := uint64(parsedData.Box.Format.Height / parsedData.Product.Thickness)

	// максимальное кол-во изделий в коробке
	maxProductQuantityInBox := maxProductQuantityInStack * impResult.Total

	// :WARNING: maxProductQuantityInBox > 0

	// кол-во полностью заполненных коробок (отбрасывается дробная часть)
	filledBoxQuantity := parsedData.Product.Quantity / maxProductQuantityInBox
	totalBoxQuantity := filledBoxQuantity

	// кол-во оставшихся изделий (последняя неполная коробка)
	restProductQuantity := parsedData.Product.Quantity - filledBoxQuantity*maxProductQuantityInBox

	// площадь изделия
	productArea := parsedData.Product.Format.Area()

	// объём коробки
	boxVolume := parsedData.Box.Format.Volume()

	// объём одного изделия
	productVolume := productArea * parsedData.Product.Thickness

	if filledBoxQuantity > 0 {
		box = entity.BoxResult{
			ProductQuantity:     maxProductQuantityInBox,
			ProductVolume:       mrlib.RoundFloat8(productVolume * float64(maxProductQuantityInBox)),
			Weight:              mrlib.RoundFloat4(productArea*float64(maxProductQuantityInBox)*parsedData.Product.WeightM2 + parsedData.Box.Weight),
			Volume:              mrlib.RoundFloat8(boxVolume),
			UnusedVolumePercent: mrlib.RoundFloat2(100 - 100*productVolume*float64(maxProductQuantityInBox)/boxVolume),
		}
	}

	if restProductQuantity > 0 {
		lastBox = entity.BoxResult{
			ProductQuantity:     restProductQuantity,
			ProductVolume:       mrlib.RoundFloat8(productVolume * float64(restProductQuantity)),
			Weight:              mrlib.RoundFloat4(productArea*float64(restProductQuantity)*parsedData.Product.WeightM2 + parsedData.Box.Weight),
			Volume:              mrlib.RoundFloat8(boxVolume),
			UnusedVolumePercent: mrlib.RoundFloat2(100 - 100*productVolume*float64(restProductQuantity)/boxVolume),
		}

		totalBoxQuantity++
	}

	uc.emitEvent(ctx, "Calc", mrmsg.Data{"raw": parsedData})

	return entity.AlgoResult{
		Box:            box,
		LastBox:        lastBox,
		ProductsVolume: mrlib.RoundFloat8(productVolume * float64(parsedData.Product.Quantity)),
		BoxesQuantity:  totalBoxQuantity,
		BoxesVolume:    mrlib.RoundFloat8(boxVolume * float64(totalBoxQuantity)),
	}, nil
}

func (uc *CirculationPackInBox) parse(data entity.RawData) (entity.ParsedData, error) {
	productFormat, err := rect.ParseFormat(data.Product.Format)
	if err != nil {
		return entity.ParsedData{}, err
	}

	boxFormat, err := parallelepiped.ParseFormat(data.Box.Format)
	if err != nil {
		return entity.ParsedData{}, err
	}

	boxMargins, err := parallelepiped.ParseFormat(data.Box.Margins)
	if err != nil {
		return entity.ParsedData{}, err
	}

	// TODO if boxMargins > boxFormat

	return entity.ParsedData{
		Product: entity.ParsedProduct{
			Format: rect.Format{
				Width:  productFormat.Width * measure.OneThousandth,  // mm -> m
				Height: productFormat.Height * measure.OneThousandth, // mm -> m
			},
			Thickness: float64(data.Product.Thickness) * measure.OneMillionth, // mkm -> m
			WeightM2:  float64(data.Product.WeightM2) * measure.OneThousandth, // g/m2 -> kg/m2
			Quantity:  data.Product.Quantity,
		},
		Box: entity.ParsedBox{
			Format: parallelepiped.Format{
				Length: boxFormat.Length * measure.OneThousandth, // mm -> m
				Width:  boxFormat.Width * measure.OneThousandth,  // mm -> m
				Height: boxFormat.Height * measure.OneThousandth, // mm -> m
			},
			Thickness: float64(data.Box.Thickness) * measure.OneMillionth, // mkm -> m
			Margins: parallelepiped.Format{
				Length: boxMargins.Length * measure.OneThousandth, // mm -> m
				Width:  boxMargins.Width * measure.OneThousandth,  // mm -> m
				Height: boxMargins.Height * measure.OneThousandth, // mm -> m
			},
			Weight:    float64(data.Box.Weight) * measure.OneThousandth,    // g -> kg
			MaxWeight: float64(data.Box.MaxWeight) * measure.OneThousandth, // g -> kg
		},
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
