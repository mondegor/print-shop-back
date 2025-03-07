package packinbox

import (
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/algo/sheet/imposition"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"
)

type (
	// Algo - размещение изделий одного формата в указанной коробке.
	Algo struct {
		logger     mrlog.Logger
		imp        *imposition.Algo
		impOptions imposition.Options
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
func (a *Algo) Calc(box model.Box, productHeap model.ProductStack) (packInBox model.PackInBox, err error) {
	if !productHeap.Format.IsValid() {
		return model.PackInBox{}, fmt.Errorf("product.Format is not valid: %s", productHeap.Format)
	}

	boxInnerFormat := box.NarrowedFormat()

	if !boxInnerFormat.IsValid() {
		return model.PackInBox{}, fmt.Errorf("box.NarrowedFormat=%s is not valid", boxInnerFormat)
	}

	if productHeap.Quantity == 0 {
		return model.PackInBox{}, nil
	}

	impResult, err := a.imp.Calc(
		productHeap.Format.BottomFormat(),
		rect2d.Format{}, // without distance
		boxInnerFormat.BottomFormat(),
		a.impOptions,
	)
	if err != nil {
		return model.PackInBox{}, fmt.Errorf(
			"packinbox.Algo.Calc[product.Format=%s, boxBottomInnerFormat=%s]: %w",
			productHeap.Format.BottomFormat(),
			boxInnerFormat.BottomFormat(),
			err,
		)
	}

	if impResult.Fragments.TotalQuantity() == 0 {
		return model.PackInBox{}, nil
	}

	// максимальное кол-во изделий в стопке (отбрасывается дробная часть)
	maxProductQuantityInStack := uint64(boxInnerFormat.Height / productHeap.Format.Height)

	if maxProductQuantityInStack == 0 {
		return model.PackInBox{}, fmt.Errorf(
			"box.NarrowedFormat.Height=%.f less than product.Height=%.f",
			boxInnerFormat.Height,
			productHeap.Format.Height,
		)
	}

	// максимальное кол-во изделий в коробке
	maxProductQuantityInBox := maxProductQuantityInStack * impResult.Fragments.TotalQuantity()

	// кол-во полностью заполненных коробок (отбрасывается дробная часть)
	filledBoxQuantity := productHeap.Quantity / maxProductQuantityInBox

	// кол-во оставшихся изделий (последняя незаполненная коробка)
	restProductQuantity := productHeap.Quantity % maxProductQuantityInBox

	if filledBoxQuantity > 0 {
		packInBox.FullBox = model.FilledBoxStack{
			FilledBox: model.FilledBox{
				Box: box,
				Product: model.ProductStack{
					Product:  productHeap.Product,
					Quantity: maxProductQuantityInBox,
				},
			},
			Quantity: filledBoxQuantity,
		}
	}

	if restProductQuantity > 0 {
		packInBox.RestBox = model.FilledBox{
			Box: box,
			Product: model.ProductStack{
				Product:  productHeap.Product,
				Quantity: restProductQuantity,
			},
		}
	}

	return packInBox, nil
}
