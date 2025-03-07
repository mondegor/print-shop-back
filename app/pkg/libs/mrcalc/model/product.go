package model

import (
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3/rect3d"
)

type (
	// Product - изделие, которое необходимо разместить в коробке.
	Product struct {
		Format rect3d.Format
		Weight float64
	}

	// ProductStack - стопка изделий.
	ProductStack struct {
		Product
		Quantity uint64
	}
)

// IsEmpty - проверяется есть ли в пачке изделия.
func (m ProductStack) IsEmpty() bool {
	return m.Quantity == 0
}

// TotalVolume - возвращает общий объём изделий.
func (m ProductStack) TotalVolume() float64 {
	return m.Format.Volume() * float64(m.Quantity)
}

// TotalWeight - возвращает общий вес продукции.
func (m ProductStack) TotalWeight() float64 {
	return m.Weight * float64(m.Quantity)
}
