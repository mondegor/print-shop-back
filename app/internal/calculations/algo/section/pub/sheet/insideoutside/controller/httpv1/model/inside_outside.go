package model

import "github.com/mondegor/print-shop-back/pkg/libs/mrcalc/s2/rect2d"

type (
	// SheetInsideOutsideQuantityRequest - comment struct.
	SheetInsideOutsideQuantityRequest struct {
		InFormat  string `json:"inFormat" validate:"required,max=16,tag_2d_size"`  // mm x mm
		OutFormat string `json:"outFormat" validate:"required,max=16,tag_2d_size"` // mm x mm
	}

	// SheetInsideOutsideMaxRequest - comment struct.
	SheetInsideOutsideMaxRequest SheetInsideOutsideQuantityRequest

	// SheetInsideOutsideQuantityResponse - результат работы алгоритма AlgoQuantity.
	SheetInsideOutsideQuantityResponse struct {
		Layout rect2d.Layout `json:"layout"`
		Total  uint64        `json:"total"`
	}

	// SheetInsideOutsideMaxResponse - результат работы алгоритма AlgoMax.
	SheetInsideOutsideMaxResponse struct {
		Fragments []rect2d.Fragment `json:"fragments"`
		Total     uint64            `json:"total"`
	}
)
