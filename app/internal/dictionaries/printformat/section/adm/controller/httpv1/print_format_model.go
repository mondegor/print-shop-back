package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/printformat/section/adm/entity"
	"github.com/mondegor/print-shop-back/pkg/libs/measure"
)

type (
	// CreatePrintFormatRequest - comment struct.
	CreatePrintFormatRequest struct {
		Caption string             `json:"caption" validate:"required,max=64"`
		Width   measure.Millimeter `json:"width" validate:"required,gte=1,lte=10000"`
		Height  measure.Millimeter `json:"height" validate:"required,gte=1,lte=10000"`
	}

	// StorePrintFormatRequest - comment struct.
	StorePrintFormatRequest struct {
		TagVersion int32              `json:"tagVersion" validate:"required,gte=1"`
		Caption    string             `json:"caption" validate:"omitempty,max=64"`
		Width      measure.Millimeter `json:"width" validate:"omitempty,gte=1,lte=10000"`
		Height     measure.Millimeter `json:"height" validate:"omitempty,gte=1,lte=10000"`
	}

	// PrintFormatListResponse - comment struct.
	PrintFormatListResponse struct {
		Items []entity.PrintFormat `json:"items"`
		Total int64                `json:"total"`
	}
)
