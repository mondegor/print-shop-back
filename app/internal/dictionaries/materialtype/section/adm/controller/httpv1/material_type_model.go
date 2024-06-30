package httpv1

import (
	"github.com/mondegor/print-shop-back/internal/dictionaries/materialtype/section/adm/entity"
)

type (
	// CreateMaterialTypeRequest - comment struct.
	CreateMaterialTypeRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	// StoreMaterialTypeRequest - comment struct.
	StoreMaterialTypeRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=64"`
	}

	// MaterialTypeListResponse - comment struct.
	MaterialTypeListResponse struct {
		Items []entity.MaterialType `json:"items"`
		Total int64                 `json:"total"`
	}
)
