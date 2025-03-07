package model

import (
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
	"github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s3/rect3d"
)

type (
	// SheetPackInStackRequest - comment struct.
	SheetPackInStackRequest struct {
		Sheet           SheetRequest `json:"sheet" validate:"required"`
		QuantityInStack uint64       `json:"quantityInStack" validate:"required,gte=1,lte=1000000"`
	}

	// SheetRequest - comment struct.
	SheetRequest struct {
		Format    string                `json:"format" validate:"required,max=16,tag_2d_size"` // mm x mm
		Thickness measure.Micrometer    `json:"thickness" validate:"required,gte=1,lte=10000"`
		Density   measure.GramPerMeter2 `json:"density" validate:"required,gte=1,lte=10000"`
		Quantity  uint64                `json:"quantity" validate:"required,gte=1,lte=1000000000"`
	}

	// SheetPackInStackResponse - результат работы алгоритма Algo.
	SheetPackInStackResponse struct {
		FullProduct   ProductResponse  `json:"fullProduct"`
		RestProduct   *ProductResponse `json:"restProduct,omitempty"`
		TotalQuantity uint64           `json:"totalQuantity"`
		TotalWeight   measure.Kilogram `json:"totalWeight"`
		TotalVolume   measure.Meter3   `json:"totalVolume"`
	}

	// ProductResponse - результаты вычислений параметров коробки.
	ProductResponse struct {
		Format rect3d.Format    `json:"format"`
		Weight measure.Kilogram `json:"weight"`
		Volume measure.Meter3   `json:"volume"`
	}
)
