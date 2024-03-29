package http_v1

import entity "print-shop-back/internal/modules/dictionaries/laminate-type/entity/admin-api"

type (
	CreateLaminateTypeRequest struct {
		Caption string `json:"caption" validate:"required,max=64"`
	}

	StoreLaminateTypeRequest struct {
		TagVersion int32  `json:"tagVersion" validate:"required,gte=1"`
		Caption    string `json:"caption" validate:"required,max=64"`
	}

	LaminateTypeListResponse struct {
		Items []entity.LaminateType `json:"items"`
		Total int64                 `json:"total"`
	}
)
