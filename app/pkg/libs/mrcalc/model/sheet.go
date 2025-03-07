package model

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"

type (
	// Sheet - изделие, которое необходимо разместить в коробке.
	Sheet struct {
		Format    rect2d.Format
		Thickness float64
		Density   float64
	}

	// SheetStack - стопка листов.
	SheetStack struct {
		Sheet
		Quantity uint64
	}
)
