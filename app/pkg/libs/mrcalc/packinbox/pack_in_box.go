package packinbox

import (
	"fmt"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/parallelepiped"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/rect/imposition"
)

type (
	// Algo - размещение изделий одного формата в указанной коробке.
	Algo struct {
		logger     mrlog.Logger
		imp        *imposition.Algo
		impOptions imposition.Options
	}

	// Product - изделие, которое необходимо разместить в коробке.
	Product struct {
		Format    rect.Format
		Thickness float64
		WeightM2  float64
		Quantity  uint64
	}

	// Box - коробка, в которой размещается изделие.
	Box struct {
		Format    parallelepiped.Format
		Thickness float64
		Margins   parallelepiped.Format
		Weight    float64
		MaxWeight float64
	}

	// AlgoResult - результат работы алгоритма Algo.
	AlgoResult struct {
		FullBox          BoxResult // полностью заполненная коробка
		RestBox          BoxResult // коробка с остатком изделий
		BoxesQuantity    uint64    // количество коробок
		BoxesWeight      float64   // общий вес коробок с изделиями
		ProductsVolume   float64   // общий объём изделий
		BoxesVolume      float64   // общий внешний объём коробок
		BoxesInnerVolume float64   // общий внутренний объём коробок
	}

	// BoxResult - результаты вычислений параметров коробки.
	BoxResult struct {
		Weight              float64 // вес коробки с изделиями
		Volume              float64 // внешний объём коробки
		InnerVolume         float64 // внутренний объём коробки
		ProductQuantity     uint64  // количество изделий в коробке
		ProductVolume       float64 // объём изделий в коробке
		UnusedVolumePercent float64 // незаполненный объём коробки в %
	}
)

// New - создаёт объект Algo.
func New(logger mrlog.Logger, imp *imposition.Algo) *Algo {
	return &Algo{
		logger: logger,
		imp:    imp,
		impOptions: imposition.Options{
			AllowRotation: true,
			UseMirror:     false,
		},
	}
}

// Calc - расчёт алгоритма.
func (a *Algo) Calc(box Box, product Product) (AlgoResult, error) {
	if product.Thickness <= 0 {
		return AlgoResult{}, fmt.Errorf("product.Thickness is not valid: %.f", product.Thickness)
	}

	const twoSides = 2
	boxSidesThickness := box.Thickness * twoSides

	// внутренний объём коробки (вычитается толщина стенок)
	boxInnerVolume := box.Format.Diff(
		parallelepiped.Format{
			Length: boxSidesThickness,
			Width:  boxSidesThickness,
			Height: boxSidesThickness,
		},
	).Volume()

	if boxInnerVolume <= 0 {
		return AlgoResult{}, fmt.Errorf(
			"box.Size=%s without boxSidesThickness=%.f is not valid: boxInnerVolume=%.f",
			box.Format,
			boxSidesThickness,
			boxInnerVolume,
		)
	}

	boxBottomFormat := box.Format.BottomFormat().Diff(box.Margins.BottomFormat())

	// вычитается толщина стенки (две стороны)
	boxBottomFormat = boxBottomFormat.Diff(
		rect.Format{
			Width:  boxSidesThickness,
			Height: boxSidesThickness,
		},
	)

	impResult, err := a.imp.Calc(
		rect.Item{
			Format: product.Format,
		},
		boxBottomFormat,
		a.impOptions,
	)
	if err != nil {
		return AlgoResult{}, fmt.Errorf("packinbox.Algo.Calc[product.Format=%s, boxBottomFormat=%s]: %w", product.Format, boxBottomFormat, err)
	}

	if impResult.Total == 0 {
		return AlgoResult{}, nil
	}

	// максимальное кол-во изделий в стопке (отбрасывается дробная часть)
	maxProductQuantityInStack := uint64((box.Format.Height - boxSidesThickness) / product.Thickness)

	if maxProductQuantityInStack == 0 {
		return AlgoResult{}, fmt.Errorf(
			"box.Format.Height=%.f without boxSidesThickness=%.f less than product.Thickness=%.f",
			box.Format.Height,
			boxSidesThickness,
			product.Thickness,
		)
	}

	// максимальное кол-во изделий в коробке
	maxProductQuantityInBox := maxProductQuantityInStack * impResult.Total

	// кол-во полностью заполненных коробок (отбрасывается дробная часть)
	filledBoxQuantity := product.Quantity / maxProductQuantityInBox
	totalBoxQuantity := filledBoxQuantity

	// кол-во оставшихся изделий (последняя незаполненная коробка)
	restProductQuantity := product.Quantity - filledBoxQuantity*maxProductQuantityInBox

	// площадь изделия
	productArea := product.Format.Area()

	// объём коробки
	boxVolume := box.Format.Volume()

	// объём одного изделия
	productVolume := productArea * product.Thickness

	var (
		fullBox     BoxResult
		restBox     BoxResult
		boxesWeight float64
	)

	if filledBoxQuantity > 0 {
		fullBox = BoxResult{
			Weight:              productArea*float64(maxProductQuantityInBox)*product.WeightM2 + box.Weight,
			Volume:              boxVolume,
			InnerVolume:         boxInnerVolume,
			ProductQuantity:     maxProductQuantityInBox,
			ProductVolume:       productVolume * float64(maxProductQuantityInBox),
			UnusedVolumePercent: 100 - 100*productVolume*float64(maxProductQuantityInBox)/boxInnerVolume,
		}

		boxesWeight += fullBox.Weight * float64(filledBoxQuantity)
	}

	if restProductQuantity > 0 {
		restBox = BoxResult{
			Weight:              productArea*float64(restProductQuantity)*product.WeightM2 + box.Weight,
			Volume:              boxVolume,
			InnerVolume:         boxInnerVolume,
			ProductQuantity:     restProductQuantity,
			ProductVolume:       productVolume * float64(restProductQuantity),
			UnusedVolumePercent: 100 - 100*productVolume*float64(restProductQuantity)/boxInnerVolume,
		}

		totalBoxQuantity++
		boxesWeight += restBox.Weight
	}

	return AlgoResult{
		FullBox:          fullBox,
		RestBox:          restBox,
		BoxesQuantity:    totalBoxQuantity,
		BoxesWeight:      boxesWeight,
		ProductsVolume:   productVolume * float64(product.Quantity),
		BoxesVolume:      boxVolume * float64(totalBoxQuantity),
		BoxesInnerVolume: boxInnerVolume * float64(totalBoxQuantity),
	}, nil
}
