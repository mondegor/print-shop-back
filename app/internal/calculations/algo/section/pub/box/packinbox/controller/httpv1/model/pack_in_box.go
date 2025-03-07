package model

import "github.com/mondegor/print-shop-back/pkg/libs/measure"

type (
	// CalcBoxPackInBoxRequest - comment struct.
	CalcBoxPackInBoxRequest struct {
		Product ProductRequest `json:"product" validate:"required"`
		Box     BoxRequest     `json:"box" validate:"required"`
	}

	// ProductRequest - comment struct.
	ProductRequest struct {
		Format   string       `json:"format" validate:"required,max=24,tag_3d_size"` // mm x mm x mm
		Weight   measure.Gram `json:"weight" validate:"required,gte=1,lte=10000"`
		Quantity uint64       `json:"quantity" validate:"required,gte=1,lte=1000000000"`
	}

	// BoxRequest - comment struct.
	BoxRequest struct {
		Format    string             `json:"format" validate:"required,max=24,tag_3d_size"` // mm x mm x mm
		Thickness measure.Micrometer `json:"thickness" validate:"required,gte=1,lte=10000"`
		Margins   string             `json:"margins" validate:"required,max=16,tag_3d_size"` // mm x mm x mm
		Weight    measure.Gram       `json:"weight" validate:"required,gte=1,lte=10000"`
		MaxWeight measure.Gram       `json:"maxWeight" validate:"omitempty,gte=1,lte=1000000"`
	}

	// BoxPackInBoxResponse - результат работы алгоритма Algo.
	BoxPackInBoxResponse struct {
		FullBox          BoxResponse      `json:"fullBox"`
		RestBox          *BoxResponse     `json:"restBox"`
		BoxesQuantity    uint64           `json:"boxesQuantity"`
		BoxesWeight      measure.Kilogram `json:"boxesWeight"`
		ProductsVolume   measure.Meter3   `json:"productsVolume"`
		BoxesVolume      measure.Meter3   `json:"boxesVolume"`
		BoxesInnerVolume measure.Meter3   `json:"boxesInnerVolume"`
	}

	// BoxResponse - результаты вычислений параметров коробки.
	BoxResponse struct {
		Weight              measure.Kilogram `json:"weight"`
		Volume              measure.Meter3   `json:"volume"`
		InnerVolume         measure.Meter3   `json:"innerVolume"`
		ProductQuantity     uint64           `json:"productQuantity"`
		ProductVolume       measure.Meter3   `json:"productVolume"`
		UnusedVolumePercent float64          `json:"unusedVolumePercent"`
	}
)
