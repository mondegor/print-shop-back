package http_v1

import (
	entity "print-shop-back/internal/modules/dictionaries/print-format/entity/admin-api"
	"print-shop-back/pkg/libs/measure"
)

type (
	CreatePrintFormatRequest struct {
		Caption string             `json:"caption" validate:"required,max=64"`
		Length  measure.Micrometer `json:"length" validate:"required,gte=1,lte=10000000"`
		Width   measure.Micrometer `json:"width" validate:"required,gte=1,lte=10000000"`
	}

	StorePrintFormatRequest struct {
		Version int32              `json:"version" validate:"required,gte=1"`
		Caption string             `json:"caption" validate:"omitempty,max=64"`
		Length  measure.Micrometer `json:"length" validate:"omitempty,gte=1,lte=10000000"`
		Width   measure.Micrometer `json:"width" validate:"omitempty,gte=1,lte=10000000"`
	}

	PrintFormatListResponse struct {
		Items []entity.PrintFormat `json:"items"`
		Total int64                `json:"total"`
	}
)
