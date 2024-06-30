package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/papercolor/section/adm/entity"
)

type (
	// CreatePaperColorRequest - comment struct.
	CreatePaperColorRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	// StorePaperColorRequest - comment struct.
	StorePaperColorRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=64"`
	}

	// PaperColorListResponse - comment struct.
	PaperColorListResponse struct {
		Items []entity.PaperColor `json:"items"`
		Total int64               `json:"total"`
	}
)
