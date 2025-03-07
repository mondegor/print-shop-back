package packinstack

import (
	"errors"
	"fmt"

	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/model"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3/rect3d"
)

type (
	// AlgoSheet - размещение изделий одного формата в указанной коробке.
	AlgoSheet struct {
		logger mrlog.Logger
	}
)

// New - создаёт объект AlgoSheet.
func New(logger mrlog.Logger) *AlgoSheet {
	return &AlgoSheet{
		logger: logger,
	}
}

// Calc - расчёт алгоритма.
func (a *AlgoSheet) Calc(sheetHeap model.SheetStack, quantityInStack uint64) (packInStack model.PackInStack, err error) {
	if !sheetHeap.Format.IsValid() {
		return model.PackInStack{}, fmt.Errorf("sheet.Format is not valid: %s", sheetHeap.Format)
	}

	if sheetHeap.Thickness <= 0 {
		return model.PackInStack{}, errors.New("sheet.Thickness must be greater than 0")
	}

	if quantityInStack == 0 {
		return model.PackInStack{}, errors.New("quantityInStack must be greater than 0")
	}

	if sheetHeap.Quantity == 0 {
		return model.PackInStack{}, nil
	}

	packQuantity := sheetHeap.Quantity / quantityInStack
	restQuantity := sheetHeap.Quantity % quantityInStack

	if packQuantity > 0 {
		packInStack.FullProduct.Product = createProduct(&sheetHeap.Sheet, quantityInStack)
		packInStack.FullProduct.Quantity = packQuantity
	}

	if restQuantity > 0 {
		packInStack.RestProduct = createProduct(&sheetHeap.Sheet, restQuantity)
	}

	return packInStack, nil
}

func createProduct(sheet *model.Sheet, quantityInStack uint64) model.Product {
	return model.Product{
		Format: rect3d.Format{
			Length: sheet.Format.Height,
			Width:  sheet.Format.Width,
			Height: sheet.Thickness * float64(quantityInStack),
		},
		Weight: sheet.Format.Area() * float64(quantityInStack) * sheet.Density,
	}
}
